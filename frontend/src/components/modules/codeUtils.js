import * as runtime from "../../../wailsjs/runtime/runtime";

// 格式化验证码为 3-3 格式
export const formatCode = (code) => {
  if (code && code.length >= 6) {
    return code.substring(0, 3) + ' ' + code.substring(3, 6);
  }
  return code;
};

// 复制验证码到剪贴板
export const copyToClipboard = async (code) => {
  try {
    const success = await runtime.ClipboardSetText(code);
    return success;
  } catch (error) {
    console.error('复制到剪贴板失败:', error);
    return false;
  }
};