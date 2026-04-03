// 1. 全局权限拦截 (只要引入了这个 JS 的页面，没登录绝对进不去)
const token = localStorage.getItem('cms_token');
if (!token) {
  alert('未登录或登录已过期，请重新登录！');
  window.location.href = '/admin/login.html';
}

// 2. 退出登录
function logout() {
  localStorage.removeItem('cms_token');
  window.location.href = '/admin/login.html';
}

// 3. 核心魔法：动态注入左侧菜单和顶部导航
// 参数 activeMenu 用来告诉菜单，当前哪个按钮应该高亮 (比如传 'news' 或 'product')
function initLayout(activeMenu) {
  // ---- 注入左侧菜单 ----
  const sidebarHtml = `
        <div class="p-3 fs-5 fw-bold text-white text-center border-bottom border-secondary" style="letter-spacing: 1px;">
            TONFY CMS
        </div>
        <div class="list-group list-group-flush mt-3 px-2">
            <a href="/admin/index.html" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'index' ? 'active bg-primary' : ''}">
                🏠 后台首页
            </a>
            <a href="/admin/news.html" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'news' ? 'active bg-primary' : ''}">
                📰 新闻管理
            </a>
            <a href="/admin/product.html" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'product' ? 'active bg-primary' : ''}">
                📦 产品管理
            </a>
            <a href="/admin/category.html" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'category' ? 'active bg-primary' : ''}">
                📂 行业分类管理
            </a>
            <a href="/admin/job.html" class="list-group-item list-group-item-action bg-dark text-white border-0 rounded mb-1 ${activeMenu === 'job' ? 'active bg-primary' : ''}">
                💼 职位管理
            </a>
        </div>
    `;
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
