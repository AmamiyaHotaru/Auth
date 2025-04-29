import { GetSecretsList, InsertSecret, DeleteSecret, RecognizeQRCode, UpdateSecret } from '../../../wailsjs/go/main/App';

// 获取账户列表
export const getSecretsList = async () => {
  try {
    const response = await GetSecretsList();
    return response;
  } catch (error) {
    console.error('获取账户列表时出错:', error);
    return [];
  }
};

// 添加新账户（手动）
export const addManualSecret = async (accountName, serverName, secret, accountType) => {
  try {
    await InsertSecret(accountName, serverName, secret, accountType);
    return true;
  } catch (error) {
    console.error('添加账户失败:', error);
    throw error;
  }
};

// 通过二维码添加账户
export const addQRCodeSecret = async (imageData) => {
  try {
    await RecognizeQRCode(imageData);
    return true;
  } catch (error) {
    console.error('二维码识别失败:', error);
    throw error;
  }
};

// 更新账户
export const updateSecret = async (id, accountName, serverName, accountType) => {
  try {
    await UpdateSecret(id, accountName, serverName, accountType);
    return true;
  } catch (error) {
    console.error('编辑账户失败:', error);
    throw error;
  }
};

// 删除账户
export const deleteSecret = async (id) => {
  try {
    await DeleteSecret(id);
    return true;
  } catch (error) {
    console.error('删除账户失败:', error);
    throw error;
  }
};

// 处理账户列表数据，添加 timeLeft 字段
export const processAccountList = (list, currentTimeLeft) => {
  if (!list || list.length === 0) {
    console.log('获取到的账户列表为空');
    return [];
  }
  
  return list.map(account => ({
    ...account, 
    timeLeft: currentTimeLeft
  }));
};