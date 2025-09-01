import Network from './network';

// 登陆相关接口
export interface LoginData {
  username: string;
  password: string;
  captchaId: string;
}

export interface LoginResponse {
  token: string;
  role_id: number;
  name: string;
  requirePasswordChange?: boolean;
}

export const login = (data: LoginData) => Network.post<LoginResponse>("/api/login", data);

// 用户CRUD操作
export const register = (data: any) => Network.post("/api/register", data);
export const getAllUsers = () => Network.get("/api/users");
export const updateUser = (id: number, data: any) => Network.put(`/api/users/${id}`, data);
export const deleteUser = (id: number) => Network.delete(`/api/users/${id}`);
export const updatePassword = (data: any) => Network.post("/api/users/update_password", data);
export const resetUserTraffic = (id: number) => Network.post(`/api/users/${id}/reset_traffic`);

// 节点CRUD操作
export const createNode = (data: any) => Network.post("/api/nodes", data);
export const getNodeList = () => Network.get("/api/nodes");
export const getNode = (id: number) => Network.get(`/api/nodes/${id}`);
export const updateNode = (id: number, data: any) => Network.put(`/api/nodes/${id}`, data);
export const deleteNode = (id: number) => Network.delete(`/api/nodes/${id}`);
export const getNodeInstallCommand = (id: number) => Network.post(`/api/nodes/${id}/install`);

// 隧道CRUD操作
export const createTunnel = (data: any) => Network.post("/api/tunnels", data);
export const getTunnelList = () => Network.get("/api/tunnels");
export const getTunnelById = (id: number) => Network.get(`/api/tunnels/${id}`);
export const updateTunnel = (id: number, data: any) => Network.put(`/api/tunnels/${id}`, data);
export const deleteTunnel = (id: number) => Network.delete(`/api/tunnels/${id}`);
export const diagnoseTunnel = (tunnelId: number) => Network.get(`/api/tunnels/${tunnelId}/diagnose`);

// 用户隧道权限管理操作
export const assignUserTunnel = (data: any) => Network.post("/api/tunnels/assign", data);
export const getUserTunnels = (id: number) => Network.get(`/api/users/${id}/tunnels`);
export const removeUserTunnel = (userId: number, tunnelId: number) => Network.delete(`/api/users/${userId}/tunnels/${tunnelId}`);

// 转发CRUD操作
export const createForward = (data: any) => Network.post("/api/forwards", data);
export const getForwardList = () => Network.get("/api/forwards");
export const getForward = (id: number) => Network.get(`/api/forwards/${id}`);
export const updateForward = (id: number, data: any) => Network.put(`/api/forwards/${id}`, data);
export const deleteForward = (id: number) => Network.delete(`/api/forwards/${id}`);
export const pauseForward = (id: number) => Network.post(`/api/forwards/${id}/pause`);
export const resumeForward = (id: number) => Network.post(`/api/forwards/${id}/resume`);
export const diagnoseForward = (id: number) => Network.get(`/api/forwards/${id}/diagnose`);
export const reorderForward = (data: any) => Network.post("/api/forwards/reorder", data);

// 转发诊断操作
// 限速规则CRUD操作
export const createSpeedLimit = (data: any) => Network.post("/api/speedlimits", data);
export const getSpeedLimitList = () => Network.get("/api/speedlimits");
export const getSpeedLimit = (id: number) => Network.get(`/api/speedlimits/${id}`);
export const updateSpeedLimit = (id: number, data: any) => Network.put(`/api/speedlimits/${id}`, data);
export const deleteSpeedLimit = (id: number) => Network.delete(`/api/speedlimits/${id}`);

// 网站配置相关接口
export const getConfigs = () => Network.get("/config/list");
export const getConfigByName = (name: string) => Network.get(`/config/get/${name}`);
export const updateConfigs = (configMap: Record<string, string>) => Network.post("/config/update", configMap);
export const updateConfig = (name: string, value: string) => Network.post("/config/update-single", { name, value });


// 验证码相关接口
export const generateCaptcha = () => Network.get("/api/captcha/generate");
export const serveCaptcha = (captchaId: string) => Network.get(`/api/captcha/${captchaId}`);
export const verifyCaptcha = (data: any) => Network.post("/api/captcha/verify", data);