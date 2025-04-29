// 进入选择模式
export const enterSelectionMode = (selectedId, selectionMode, selectedAccountIds) => {
  selectionMode.value = true;
  if (selectedId) {
    selectedAccountIds.value = new Set([selectedId]);
  }
};

// 退出选择模式
export const cancelSelectionMode = (selectionMode, selectedAccountIds) => {
  selectionMode.value = false;
  selectedAccountIds.value.clear();
};

// 切换账户的选择状态
export const toggleSelection = (id, selectedAccountIds) => {
  if (selectedAccountIds.value.has(id)) {
    selectedAccountIds.value.delete(id);
  } else {
    selectedAccountIds.value.add(id);
  }
};

// 请求删除选中的账户
export const requestDeleteSelected = (selectedCount, selectedAccountIds) => {
  if (selectedCount.value === 0) return;
  
  let title = '删除确认';
  let message = `确定要删除这${selectedCount.value}个账户吗？此操作不可撤销。`;
  
  return {
    title,
    message
  };
};