/**
 * create by zhangxiang on 2023-06-06 10:23
 * 类注释：
 * 备注：
 */
const req = require.context('./svg', true, /\.svg$/)
const requireAll = requireContext => requireContext.keys().map(requireContext)
requireAll(req)



