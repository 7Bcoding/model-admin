// 特殊Feature配置定义
// 这些feature会作为key-value对存储在feature_override中

// SupportedFeatures should be defined as a constant array
export const SupportedEdnpointFeatures = ["/chat/completions", "/completions", "/v1/messages"];
export const SupportedResponseFormatFeatures = ["json_schema", "json_object", "text"];

// 字符集编码选项
export const InputTokenEncodingOptions = ["deepseek"];

// 特殊Feature键定义
export const FeatureKeys = {
  TTFT_TIMEOUT: 'ttft_timeout',
  TPOT_TIMEOUT: 'tpot_timeout',
  HANDLE_SEPARATED_REASONING: 'handle_separated_reasoning',
  HANDLE_COMBINED_REASONING: 'handle_combined_reasoning',
  REMOVE_CONTENT_FILTER: 'remove_content_filter',
  ADD_CONTINUE_CHUNK_USAGE: 'add_continue_chunk_usage',
  REMOVE_CONTINUOUS_USAGE_STATS: 'remove_continuous_usage_stats',
  LIMIT_MAX_TOKEN: 'limit_max_token',
  HANDLE_CONTENT_FILTER: 'handle_content_filter'
}

// 特殊Feature描述信息
export const FeatureDescriptions = {
  [FeatureKeys.TTFT_TIMEOUT]: '设置后在stream模式最多等待此时间，示例：1000ms/10s',
  [FeatureKeys.TPOT_TIMEOUT]: 'streaming模式下单chunk等待时间超时后会自动结束请求，示例：1000ms/10s',
  [FeatureKeys.HANDLE_SEPARATED_REASONING]: '当provider的reasoning内容单独返回时候，需要加上此标记',
  [FeatureKeys.HANDLE_COMBINED_REASONING]: '当provider的reasoning内容和content内容一起返回时候，需要加上此标记',
  [FeatureKeys.REMOVE_CONTENT_FILTER]: '选择此选项时，会在请求中增加额外的字段以关闭内容审核',
  [FeatureKeys.ADD_CONTINUE_CHUNK_USAGE]: '额外使用chunk_include_usage字段来获取每个chunk返回usage的能力',
  [FeatureKeys.REMOVE_CONTINUOUS_USAGE_STATS]: '移除continuous_usage_stats字段',
  [FeatureKeys.LIMIT_MAX_TOKEN]: '当传输的max_tokens大于此值时，会修正为该值，常用于reasoning-content不单独限制的推理模型',
  [FeatureKeys.HANDLE_CONTENT_FILTER]: '当遇到内容过滤的情形时，会先报错以触发重试'
}

export const FeatureName = {
  [FeatureKeys.TTFT_TIMEOUT]: 'TTFT 超时时间',
  [FeatureKeys.TPOT_TIMEOUT]: 'TPOT 超时时间',
  [FeatureKeys.HANDLE_SEPARATED_REASONING]: '[特征]reasoning内容为单独返回',
  [FeatureKeys.HANDLE_COMBINED_REASONING]: '[特征]reasoning内容和content内容一起返回',
  [FeatureKeys.REMOVE_CONTENT_FILTER]: '关闭内容审核(百度专用)',
  [FeatureKeys.ADD_CONTINUE_CHUNK_USAGE]: '使用chunk_include_usage(火山)',
  [FeatureKeys.REMOVE_CONTINUOUS_USAGE_STATS]: '移除continuous_usage_stats字段(Fireworks)',
  [FeatureKeys.LIMIT_MAX_TOKEN]: '修正max_tokens(火山)',
  [FeatureKeys.HANDLE_CONTENT_FILTER]: '内容过滤重试(火山)'
}

// 增加特殊key的类型定义，最终的格式都是字符串，但是某些情况下我们为了显示方便需要做一些类型的预定义
export const FeatureType = {
  [FeatureKeys.HANDLE_SEPARATED_REASONING]: 'boolean',
  [FeatureKeys.HANDLE_COMBINED_REASONING]: 'boolean',
  [FeatureKeys.REMOVE_CONTENT_FILTER]: 'boolean',
  [FeatureKeys.ADD_CONTINUE_CHUNK_USAGE]: 'boolean',
  [FeatureKeys.REMOVE_CONTINUOUS_USAGE_STATS]: 'boolean',
  [FeatureKeys.LIMIT_MAX_TOKEN]: 'integer',
  [FeatureKeys.HANDLE_CONTENT_FILTER]: 'boolean'
}

// 特殊Feature配置项列表
export const SpecialFeatures = [
  {
    key: FeatureKeys.TTFT_TIMEOUT,
    label: FeatureName[FeatureKeys.TTFT_TIMEOUT],
    description: FeatureDescriptions[FeatureKeys.TTFT_TIMEOUT]
  },
  {
    key: FeatureKeys.TPOT_TIMEOUT,
    label: FeatureName[FeatureKeys.TPOT_TIMEOUT], 
    description: FeatureDescriptions[FeatureKeys.TPOT_TIMEOUT]
  },
  {
    key: FeatureKeys.HANDLE_SEPARATED_REASONING,
    label: FeatureName[FeatureKeys.HANDLE_SEPARATED_REASONING],
    description: FeatureDescriptions[FeatureKeys.HANDLE_SEPARATED_REASONING]
  },
  {
    key: FeatureKeys.HANDLE_COMBINED_REASONING,
    label: FeatureName[FeatureKeys.HANDLE_COMBINED_REASONING],
    description: FeatureDescriptions[FeatureKeys.HANDLE_COMBINED_REASONING]
  },
  {
    key: FeatureKeys.REMOVE_CONTENT_FILTER,
    label: FeatureName[FeatureKeys.REMOVE_CONTENT_FILTER],
    description: FeatureDescriptions[FeatureKeys.REMOVE_CONTENT_FILTER]
  },
  {
    key: FeatureKeys.ADD_CONTINUE_CHUNK_USAGE,
    label: FeatureName[FeatureKeys.ADD_CONTINUE_CHUNK_USAGE],
    description: FeatureDescriptions[FeatureKeys.ADD_CONTINUE_CHUNK_USAGE]
  },
  {
    key: FeatureKeys.REMOVE_CONTINUOUS_USAGE_STATS,
    label: FeatureName[FeatureKeys.REMOVE_CONTINUOUS_USAGE_STATS],
    description: FeatureDescriptions[FeatureKeys.REMOVE_CONTINUOUS_USAGE_STATS]
  },
  {
    key: FeatureKeys.LIMIT_MAX_TOKEN,
    label: FeatureName[FeatureKeys.LIMIT_MAX_TOKEN],
    description: FeatureDescriptions[FeatureKeys.LIMIT_MAX_TOKEN]
  },
  {
    key: FeatureKeys.HANDLE_CONTENT_FILTER,
    label: FeatureName[FeatureKeys.HANDLE_CONTENT_FILTER],
    description: FeatureDescriptions[FeatureKeys.HANDLE_CONTENT_FILTER]
  }
]

// 获取特殊Feature的描述
export const getFeatureDescription = (key) => {
  return FeatureDescriptions[key] || ''
}

// 验证特殊Feature的值
export const validateFeatureValue = (key, value) => {
  if (!value || value.trim() === '') {
    return true // 空值是允许的
  }
  
  // 对于超时时间类型的feature，验证是否为正数
  if (key === FeatureKeys.TTFT_TIMEOUT || key === FeatureKeys.TPOT_TIMEOUT) {
    const num = parseFloat(value)
    return !isNaN(num) && num > 0
  }
  
  return true
}

// 格式化特殊Feature的值
export const formatFeatureValue = (key, value) => {
  if (!value || value.trim() === '') {
    return ''
  }
  
  // 对于超时时间类型的feature，确保是数字格式
  if (key === FeatureKeys.TTFT_TIMEOUT || key === FeatureKeys.TPOT_TIMEOUT) {
    const num = parseFloat(value)
    return isNaN(num) ? '' : num.toString()
  }
  
  return value.toString()
}

// 支付方式枚举定义
export const PaymentMethod = {
  UNKNOWN: 0,    // 未知
  PREPAID: 1,    // 预付费
  POSTPAID: 2    // 后付费
}

// 支付方式标签
export const PaymentMethodLabels = {
  [PaymentMethod.UNKNOWN]: '未知',
  [PaymentMethod.PREPAID]: '预付费',
  [PaymentMethod.POSTPAID]: '后付费'
}

// 支付方式标签类型
export const PaymentMethodTypes = {
  [PaymentMethod.UNKNOWN]: 'info',
  [PaymentMethod.PREPAID]: 'success',
  [PaymentMethod.POSTPAID]: 'warning'
}