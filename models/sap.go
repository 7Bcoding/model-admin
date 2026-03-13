package models

import serverlessv1beta1 "llm-ops/serverless/v1beta1"

// build a object to get parama from body
type SapParam struct {
	MinReplicas          int32 `json:"minReplicas"`          // 最小副本数
	MaxReplicas          int32 `json:"maxReplicas"`          // 最大副本数
	ConcurrencyPerWorker int   `json:"concurrencyPerWorker"` // 触发扩容的并发数
	ScaleUpWindow        int32 `json:"scaleUpWindow"`        // 扩容窗口时间
	ScaleDownWindow      int32 `json:"scaleDownWindow"`      // 缩容窗口时间
}

func SapParamFromSap(sap *serverlessv1beta1.ServerlessAutoScalingPolicy) *SapParam {
	ret := &SapParam{
		MinReplicas:     *sap.Spec.MinReplicas,
		MaxReplicas:     sap.Spec.MaxReplicas,
		ScaleUpWindow:   *sap.Spec.Behavior.ScaleUp.StabilizationWindowSeconds,
		ScaleDownWindow: *sap.Spec.Behavior.ScaleDown.StabilizationWindowSeconds,
	}

	if sap.Spec.Metrics != nil {
		ret.ConcurrencyPerWorker = int(*sap.Spec.Metrics[0].Resource.Target.AverageValueAsInt)
	}
	return ret
}
