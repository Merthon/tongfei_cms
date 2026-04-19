// ==========================================
// 1. 全局登录拦截
// ==========================================
if (window.location.pathname.indexOf('login.html') === -1) {
  const token = localStorage.getItem('cms_token');
  if (!token) {
    alert('未登录或登录已过期，请重新登录！');
    window.location.href = '/admin/login.html';
  }
}

// ==========================================
// 2. 核心魔法：全局统一拦截 403 权限错误 (保留这个作为最后一道防线)
// ==========================================
const originalFetch = window.fetch;
window.fetch = async function (...args) {
  const response = await originalFetch(...args);
  if (response.status === 403) {
    alert('⛔ 越权警告：抱歉，您当前账号没有操作此模块的权限！');
  }
  return response;
};

// ==========================================
// 3. 退出登录
// ==========================================
function logout() {
  localStorage.removeItem('cms_token');
  localStorage.removeItem('cms_role');
  localStorage.removeItem('cms_modules');
  window.location.href = '/admin/login.html';
}

// ==========================================
// 🚨 新增：前端点击菜单时的“前置拦截器”
// ==========================================
function tryNavigate(requiredModule, targetUrl) {
  const role = localStorage.getItem('cms_role');
  const modules = localStorage.getItem('cms_modules') || '';

  // 1. 首页不需要权限，任何人都可以进
  if (requiredModule === 'index') {
    window.location.href = targetUrl;
    return;
  }

  // 2. 超级管理员，无敌放行，去哪都行
  if (role === 'super_admin') {
    window.location.href = targetUrl;
    return;
  }

  // 3. 如果点击的是“账号管理”这种只有超管能进的页面
  if (requiredModule === 'super_admin') {
    alert('⛔ 越权拦截：只有超级管理员才能进入账号管理中心！');
    return;
  }

  // 4. 普通管理员，检查他的权限字符串里有没有这个模块
  if (modules.includes(requiredModule)) {
    window.location.href = targetUrl; // 有权限，放行跳转
  } else {
    alert('⛔ 您没有该模块的管理权限，请联系超级管理员开通！'); // 没权限，直接弹窗，页面原地不动
  }
}
// ==========================================
// 4. 动态注入左侧菜单和顶部导航
// ==========================================
function initLayout(activeMenu) {
    // 🚨 核心逻辑：从缓存中读取当前登录人的身份
    const role = localStorage.getItem('cms_role');

    // ---- 1. 拼接所有人都能看见的基础菜单 ----
    let sidebarHtml = `
        <div class="p-3 fs-5 fw-bold text-white text-center border-bottom border-secondary" style="letter-spacing: 1px;">
            TONFY CMS
        </div>
        <div class="list-group list-group-flush mt-3 px-2">
            <a href="javascript:void(0)" onclick="tryNavigate('index', '/admin/index.html')" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'index' ? 'active bg-primary' : ''}">
                🏠 后台首页
            </a>
            <a href="javascript:void(0)" onclick="tryNavigate('news', '/admin/news.html')" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'news' ? 'active bg-primary' : ''}">
                📰 新闻管理
            </a>
            <a href="javascript:void(0)" onclick="tryNavigate('product', '/admin/product.html')" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'product' ? 'active bg-primary' : ''}">
                📦 产品管理
            </a>
            <a href="javascript:void(0)" onclick="tryNavigate('product', '/admin/category.html')" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'category' ? 'active bg-primary' : ''}">
                📂 行业分类管理
            </a>
            <a href="javascript:void(0)" onclick="tryNavigate('job', '/admin/job.html')" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'job' ? 'active bg-primary' : ''}">
                💼 职位管理
            </a>
            <a href="javascript:void(0)" onclick="tryNavigate('banner', '/admin/banner.html')" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'banner' ? 'active bg-primary' : ''}">
                📺 首页banner管理
            </a>
    `;

    // ---- 2. 🚨 判断身份：只有超级管理员，才拼接账号管理菜单 ----
    if (role === 'super_admin') {
        sidebarHtml += `
            <a href="javascript:void(0)" onclick="tryNavigate('super_admin', '/admin/account.html')" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mt-3 ${activeMenu === 'account' ? 'active bg-danger' : ''}">
                🛡️ 账号管理 (超管专属)
            </a>
        `;
    }

    // ---- 3. 收尾闭合标签 ----
    sidebarHtml += `</div>`;

    const sidebarContainer = document.getElementById('sidebar-container');
    if (sidebarContainer) sidebarContainer.innerHTML = sidebarHtml;
    
    // ---- 注入顶部导航 ----
    const headerHtml = `
        <nav class="navbar navbar-light bg-white shadow-sm px-4">
            <span class="navbar-brand mb-0 h1 fs-5 text-muted">后台管理中心</span>
            <button class="btn btn-outline-danger btn-sm" onclick="logout()">退出登录</button>
        </nav>
    `;
    const headerContainer = document.getElementById('header-container');
    if (headerContainer) headerContainer.innerHTML = headerHtml;
}
