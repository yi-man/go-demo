package cluster

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// space resource limit
var (
	// 本地临时存储容量限制
	resourceEphemeralStorageLimit = resource.MustParse("50Gi") // Ki | Mi | Gi
)

// pvc volumes
var (
	// shm 内存类型，不设置大小
	shmVolume = v1.Volume{
		Name: "shm-volume",
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{
				Medium: v1.StorageMediumMemory,
			},
		},
	}
)

func GetClaimNameByRegion(prefix string, region string) string {
	// region命名规范：ksc-gic-worker, ksc-tr-worker
	fields := strings.Split(region, "-")
	if len(fields) < 3 {
		return ""
	}
	return fmt.Sprintf("%s-%s", prefix, fields[1])
}

func (cc *ClusterClient) ListPods() ([]v1.Pod, error) {
	podList, err := cc.KubeClient.CoreV1().Pods(cc.NameSpace).List(
		context.Background(),
		metav1.ListOptions{},
	)

	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}
func (cc *ClusterClient) GetPodsByJobName(ctx context.Context, jobName string) ([]v1.Pod, error) {
	podList, err := cc.KubeClient.CoreV1().Pods(cc.NameSpace).List(
		ctx,
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", jobName),
		},
	)

	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}

func (cc *ClusterClient) GetUserEvents(ctx context.Context, name string) ([]string, error) {
	events, err := cc.GetPodEvents(ctx, name)
	if err != nil {
		return nil, err
	}
	var userEventsDesc []string
	for _, event := range events {
		switch event.Reason {
		case "Pulling":
			// 正在拉取镜像
			userEventsDesc = append(userEventsDesc, "Pulling image...")
		case "ImagePullBackOff", "ErrImagePull":
			// 镜像拉取相关事件
			userEventsDesc = append(userEventsDesc, "Failed to pull image")
		case "Failed", "FailedScheduling", "FailedBinding":
			// Pod 创建失败相关事件
			userEventsDesc = append(userEventsDesc, "Failed to create pod")
		case "Unhealthy":
			// Pod 运行失败相关事件
			userEventsDesc = append(userEventsDesc, "Pod is not ready")
		case "OOMKilling":
			// Pod OOM 相关事件
			userEventsDesc = append(userEventsDesc, "Out of memory error")
		}
	}
	return userEventsDesc, nil
}

func (cc *ClusterClient) GetPodEvents(ctx context.Context, name string) ([]v1.Event, error) {
	events, err := cc.KubeClient.CoreV1().Events(cc.NameSpace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
	})
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}

func (cc *ClusterClient) GetPodRecentEventReason(ctx context.Context, name string) (reason string, err error) {
	events, err := cc.KubeClient.CoreV1().Events(cc.NameSpace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
	})
	if err != nil {
		return "", err
	}
	if len(events.Items) == 0 {
		return "", nil
	}

	for _, event := range events.Items {
		if event.Source.Component != "kubelet" && event.Type == "Normal" {
			continue
		}

		reason = event.Reason
	}
	return reason, nil
}

func (cc *ClusterClient) DeletePod(ctx context.Context, podName string) error {
	err := cc.KubeClient.CoreV1().Pods(cc.NameSpace).Delete(
		ctx,
		podName,
		metav1.DeleteOptions{
			PropagationPolicy:  &deletePolicy,
			GracePeriodSeconds: &gracePeriodSeconds,
		},
	)
	return err
}

func ValidateAndExtractContainerName(podName string) (bool, string) {
	parts := strings.Split(podName, "-")
	// 确保Pod名称至少由四部分组成，以满足[prefix]-[containerName]-[additionalInfo]的格式
	if len(parts) < 4 {
		return false, ""
	}
	containerName := strings.Join(parts[1:3], "-")
	// 这里可以添加额外的校验逻辑，比如校验容器名称的特定格式
	return true, containerName
}

// GetPodLog 根据Pod名称获取实时日志流
// region 集群名称
// name Pod名称
// tailLines 获取历史日志的行数，最多200行，对大于200KB的日志进行截断
func (cc *ClusterClient) GetPodLog(ctx context.Context, name string, tailLines int64) (io.ReadCloser, error) {
	var (
		defaultLogLines int64 = 200
		maxLogBytes     int64 = 1024 * 2000 // 2MB
	)
	if tailLines == 0 || tailLines > 200 {
		tailLines = defaultLogLines
	}

	// ksc会起一个sidecar容器，需要指明主容器
	valid, containerName := ValidateAndExtractContainerName(name)
	var opts *v1.PodLogOptions
	if valid {
		logx.Infof("containerName: %s", containerName)
		opts = &v1.PodLogOptions{
			Container:  containerName,
			TailLines:  &tailLines,
			LimitBytes: &maxLogBytes,
			Follow:     true,
		}
	} else {
		logx.Infof("containerName not valid")
		opts = &v1.PodLogOptions{
			TailLines:  &tailLines,
			LimitBytes: &maxLogBytes,
			Follow:     true,
		}
	}

	req := cc.KubeClient.CoreV1().Pods(cc.NameSpace).GetLogs(name, opts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		logx.Error("stream error: ", err.Error())
		return nil, err
	}

	return podLogs, nil
}

func (cc *ClusterClient) CreatePrivateImagePodSpec(name string, image string, args []string, region string, pvcName string) *v1.PodSpec {
	podSpec := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  name,
				Image: image,
				Command: []string{
					"/bin/sh",
					"-c",
					// "sleep 86400",
				},
				Args: args,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-grpc",
						ContainerPort: 8001,
						Protocol:      v1.ProtocolTCP,
					},
					{
						Name:          "port-http",
						ContainerPort: 8000,
						Protocol:      v1.ProtocolTCP,
					},
					{
						Name:          "port-monitor",
						ContainerPort: 8002,
						Protocol:      v1.ProtocolTCP,
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/models_storage/ckpt/custom",
						SubPath:   "common/custom_ckpt",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/models_storage/lora/custom",
						SubPath:   "common/custom_loras",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/models_storage/embedding/custom",
						SubPath:   "common/custom_embeddings",
						ReadOnly:  true,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(8000),
							Path: "/v2/health/ready",
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
			},
		},
		Volumes: []v1.Volume{
			shmVolume,
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
			{
				Name: "custom-data",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			},
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
	}
	return podSpec
}

func (cc *ClusterClient) CreateImagePodSpec(name string, image string, args []string, region string) *v1.PodSpec {
	podSpec := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  name,
				Image: image,
				Command: []string{
					"/bin/sh",
					"-c",
					// "sleep 86400",
				},
				Args: args,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-grpc",
						ContainerPort: 8001,
						Protocol:      v1.ProtocolTCP,
					},
					{
						Name:          "port-http",
						ContainerPort: 8000,
						Protocol:      v1.ProtocolTCP,
					},
					{
						Name:          "port-monitor",
						ContainerPort: 8002,
						Protocol:      v1.ProtocolTCP,
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(8000),
							Path: "/v2/health/ready",
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
			},
		},
		Volumes: []v1.Volume{
			shmVolume,
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
	}
	return podSpec
}

func (cc *ClusterClient) CreateTestPodSpec(name string, image string, region string, pvcName string) *v1.PodSpec {
	podSpec := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  name,
				Image: image,
				Command: []string{
					"/bin/sh",
					"-c",
					"sleep 864000000",
				},
				// Args: args,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-grpc",
						ContainerPort: 8001,
						Protocol:      v1.ProtocolTCP,
					},
					{
						Name:          "port-http",
						ContainerPort: 8000,
						Protocol:      v1.ProtocolTCP,
					},
					{
						Name:          "port-monitor",
						ContainerPort: 8002,
						Protocol:      v1.ProtocolTCP,
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/models_storage/ckpt/custom",
						SubPath:   "common/custom_ckpt",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/models_storage/lora/custom",
						SubPath:   "common/custom_loras",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/models_storage/embedding/custom",
						SubPath:   "common/custom_embeddings",
						ReadOnly:  true,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(8000),
							Path: "/v2/health/ready",
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
			},
		},
		Volumes: []v1.Volume{
			shmVolume,
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
			{
				Name: "custom-data",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			},
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
	}
	return podSpec
}

func (cc *ClusterClient) CreateGeneralPodSpec(name string, image string, args []string, region string) *v1.PodSpec {
	podSpec := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  name,
				Image: image,
				Args:  args,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-http",
						ContainerPort: 8000,
						Protocol:      v1.ProtocolTCP,
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/models",
						SubPath:   "models-all/models_repo",
						ReadOnly:  true,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						HTTPGet: &v1.HTTPGetAction{
							Port: intstr.FromInt(8000),
							Path: "/health",
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
			},
		},
		Volumes: []v1.Volume{
			shmVolume,
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
	}
	return podSpec
}

func (cc *ClusterClient) CreateSdApiPodSpec(containerName, image string, region string) *v1.PodSpec {
	pod := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  containerName,
				Image: image,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-1",
						ContainerPort: 7861,
						Protocol:      v1.ProtocolTCP,
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.FromInt(7861),
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			}, shmVolume,
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
	}
	return pod
}

func (cc *ClusterClient) CreatePrivateSdApiPodSpec(containerName, image string, region string, pvcName string) *v1.PodSpec {
	pod := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  containerName,
				Image: image,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-1",
						ContainerPort: 7861,
						Protocol:      v1.ProtocolTCP,
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/stable-diffusion-webui/models/Stable-diffusion/custom",
						SubPath:   "common/custom_ckpt",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/stable-diffusion-webui/models/Lora/custom",
						SubPath:   "common/custom_loras",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/stable-diffusion-webui/embeddings/custom",
						SubPath:   "common/custom_embeddings",
						ReadOnly:  true,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.FromInt(7861),
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
			{
				Name: "custom-data",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			},
			shmVolume,
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
	}
	return pod
}

func (cc *ClusterClient) CreateComfyApiPodSpec(containerName, image string, region string, creatorId int64) *v1.PodSpec {
	pod := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  containerName,
				Image: image,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-1",
						ContainerPort: 8188,
						Protocol:      v1.ProtocolTCP,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.FromInt(8188),
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       5,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
					{
						Name:      "comfyapi-common",
						MountPath: "/opt/ComfyUI/output",
						SubPath:   fmt.Sprintf("%d/output", creatorId),
						ReadOnly:  false,
					},
					{
						Name:      "comfyapi-common",
						MountPath: "/opt/ComfyUI/input",
						SubPath:   fmt.Sprintf("%d/input", creatorId),
						ReadOnly:  false,
					},
				},
			},
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
			{
				Name: "comfyapi-common",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("comfyapi", region),
					},
				},
			},
			shmVolume,
		},
	}
	return pod
}

func (cc *ClusterClient) CreatePrivateComfyApiPodSpec(containerName, image string, region string, pvcName string, serviceId int64) *v1.PodSpec {
	pod := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  containerName,
				Image: image,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-1",
						ContainerPort: 8188,
						Protocol:      v1.ProtocolTCP,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.FromInt(8188),
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       5,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/ComfyUI/models/checkpoints/custom",
						SubPath:   "common/custom_ckpt",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/ComfyUI/models/loras/custom",
						SubPath:   "common/custom_loras",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/ComfyUI/models/embeddings/custom",
						SubPath:   "common/custom_embeddings",
						ReadOnly:  true,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/ComfyUI/output",
						SubPath:   fmt.Sprintf("endpoints/%d/output", serviceId),
						ReadOnly:  false,
					},
					{
						Name:      "custom-data",
						MountPath: "/opt/ComfyUI/input",
						SubPath:   fmt.Sprintf("endpoints/%d/input", serviceId),
						ReadOnly:  false,
					},
				},
			},
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
		Volumes: []v1.Volume{
			{
				Name: "custom-data",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			},
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
			shmVolume,
		},
	}
	return pod
}

func (cc *ClusterClient) CreateCustomPodSpec(containerName, image string, region string) *v1.PodSpec {
	pod := &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  containerName,
				Image: image,
				Ports: []v1.ContainerPort{
					{
						Name:          "port-1",
						ContainerPort: 8000,
						Protocol:      v1.ProtocolTCP,
					},
				},
				ReadinessProbe: &v1.Probe{
					ProbeHandler: v1.ProbeHandler{
						TCPSocket: &v1.TCPSocketAction{
							Port: intstr.FromInt(8000),
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Lifecycle: &v1.Lifecycle{
					PreStop: &v1.LifecycleHandler{
						Exec: &v1.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep 15",
							},
						},
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "shm-volume",
						MountPath: "/dev/shm",
					},
					{
						Name:      "ai-models",
						MountPath: "/models_storage/stable-diffusion-models",
						SubPath:   "models-all/stable-diffusion-models",
						ReadOnly:  true,
					},
				},
			},
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: cc.HarborSecretKey,
			},
		},
		Volumes: []v1.Volume{
			shmVolume,
			{
				Name: "ai-models",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: GetClaimNameByRegion("ai-models", region),
					},
				},
			},
		},
	}
	return pod
}
