// 计算当前时间在 30 秒周期内的剩余秒数
export const calculateTimeLeft = () => {
  const nowSeconds = Math.floor(Date.now() / 1000);
  return 30 - (nowSeconds % 30);
};

// 计算进度条偏移量
export const calculateProgressOffset = (timeLeft) => {
  const radius = 15;
  const circumference = 2 * Math.PI * radius;
  const percentage = timeLeft / 30;
  return circumference * (1 - percentage);
};

// 检测是否需要刷新验证码
export const needsCodeRefresh = (accounts, newTimeLeft) => {
  if (accounts.length > 0) {
    return newTimeLeft === 30 && accounts[0].timeLeft === 1;
  }
  return false;
};

// 更新账户列表中的所有 timeLeft 值
export const updateAccountsTimeLeft = (accounts, newTimeLeft) => {
  return accounts.map(acc => ({
    ...acc,
    timeLeft: newTimeLeft
  }));
};

// 计算圆形进度条的样式
export const getProgressStyle = (timeLeft) => {
  // 计算百分比，确保不低于 0
  const percentage = Math.max(0, timeLeft / 30) * 100;
  return {
    background: `radial-gradient(white 60%, transparent 61%), conic-gradient(#4285F4 ${percentage}%, #e0e0e0 ${percentage}%)`
  };
};