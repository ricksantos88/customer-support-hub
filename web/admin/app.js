// ── Global Variables & Configurations ──
let currentAdminId = null;
let currentAdminEmail = null;
let currentAdminName = null;

let agentsLimit = 10;
let agentsOffset = 0;

// ── Application Initialization ──
document.addEventListener("DOMContentLoaded", () => {
    setupEventListeners();
    checkAuthentication();
});

// ── Event Listeners Binding ──
function setupEventListeners() {
    // Login Form Submit
    document.getElementById("login-form").addEventListener("submit", handleLogin);

    // Logout Button
    document.getElementById("logout-btn").addEventListener("click", handleLogout);

    // Modal Triggers
    document.getElementById("btn-open-create").addEventListener("click", () => openModal("modal-create"));

    // Forms inside Modals
    document.getElementById("create-agent-form").addEventListener("submit", handleCreateAgent);
    document.getElementById("edit-agent-form").addEventListener("submit", handleUpdateAgent);

    // Pagination
    document.getElementById("btn-prev-page").addEventListener("click", () => {
        if (agentsOffset >= agentsLimit) {
            agentsOffset -= agentsLimit;
            loadAgents();
        }
    });
    document.getElementById("btn-next-page").addEventListener("click", () => {
        agentsOffset += agentsLimit;
        loadAgents();
    });
}

// ── Authentication Checks ──
function checkAuthentication() {
    const token = localStorage.getItem("access_token");
    const role = localStorage.getItem("agent_role");

    if (token && role === "admin") {
        currentAdminId = localStorage.getItem("agent_id");
        currentAdminEmail = localStorage.getItem("agent_email");
        currentAdminName = localStorage.getItem("agent_name") || "Administrador";

        document.getElementById("user-info").innerText = `${currentAdminName} (${currentAdminEmail})`;
        showDashboard();
    } else {
        showLogin();
    }
}

function showLogin() {
    document.getElementById("login-area").classList.remove("hidden");
    document.getElementById("dashboard-area").classList.add("hidden");
}

function showDashboard() {
    document.getElementById("login-area").classList.add("hidden");
    document.getElementById("dashboard-area").classList.remove("hidden");

    // Load metrics, agents, sessions
    loadDiagnostics();
    loadAgents();
    loadSessions();

    // Poll diagnostics and sessions every 15 seconds
    setInterval(loadDiagnostics, 15000);
    setInterval(loadSessions, 15000);
}

// ── AJAX / Request Wrapper with Auth Injection and Refresh Token rotation ──
async function apiRequest(url, options = {}) {
    let token = localStorage.getItem("access_token");
    if (!options.headers) {
        options.headers = {};
    }
    if (token) {
        options.headers["Authorization"] = `Bearer ${token}`;
    }
    if (!options.headers["Content-Type"]) {
        options.headers["Content-Type"] = "application/json";
    }

    try {
        let response = await fetch(url, options);

        if (response.status === 401) {
            // Attempt to refresh token
            const refreshed = await attemptTokenRefresh();
            if (refreshed) {
                token = localStorage.getItem("access_token");
                options.headers["Authorization"] = `Bearer ${token}`;
                response = await fetch(url, options);
            } else {
                showToast("Sessão expirada. Faça login novamente.", "error");
                performLocalLogout();
                return null;
            }
        }
        return response;
    } catch (err) {
        slogError("API request error", err);
        showToast("Erro de conexão com o servidor.", "error");
        return null;
    }
}

async function attemptTokenRefresh() {
    const refreshToken = localStorage.getItem("refresh_token");
    if (!refreshToken) return false;

    try {
        const response = await fetch("/auth/refresh", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ refresh_token: refreshToken })
        });

        if (response.ok) {
            const data = await response.json();
            localStorage.setItem("access_token", data.access_token);
            localStorage.setItem("refresh_token", data.refresh_token);
            return true;
        }
    } catch (e) {
        slogError("token refresh failed", e);
    }
    return false;
}

// ── Handlers ──
async function handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;
    const loginBtn = document.getElementById("login-btn");

    loginBtn.innerText = "Entrando...";
    loginBtn.disabled = true;

    try {
        const response = await fetch("/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            const data = await response.json();

            // Decode payload from token to check role (to be safe)
            const tokenParts = data.access_token.split('.');
            if (tokenParts.length === 3) {
                const payload = JSON.parse(atob(tokenParts[1]));
                if (payload.agent_role !== "admin") {
                    showToast("Acesso negado. Apenas administradores podem entrar.", "error");
                    loginBtn.innerText = "Entrar";
                    loginBtn.disabled = false;
                    return;
                }

                localStorage.setItem("access_token", data.access_token);
                localStorage.setItem("refresh_token", data.refresh_token);
                localStorage.setItem("agent_id", data.agent_id);
                localStorage.setItem("agent_role", payload.agent_role);
                localStorage.setItem("agent_email", email);

                // Try to find the agent name from payload or default to email
                localStorage.setItem("agent_name", payload.name || "Administrador");

                showToast("Autenticado com sucesso!", "success");
                checkAuthentication();
            }
        } else {
            const data = await response.json().catch(() => ({}));
            showToast(data.error || "E-mail ou senha incorretos.", "error");
        }
    } catch (err) {
        showToast("Erro ao conectar ao servidor.", "error");
    } finally {
        loginBtn.innerText = "Entrar";
        loginBtn.disabled = false;
    }
}

async function handleLogout() {
    const response = await apiRequest("/auth/logout", { method: "POST" });
    if (response && response.ok) {
        showToast("Você saiu com sucesso.", "success");
    }
    performLocalLogout();
}

function performLocalLogout() {
    localStorage.clear();
    currentAdminId = null;
    currentAdminEmail = null;
    currentAdminName = null;
    showLogin();
}

// ── CRUD Operations ──

// List Diagnostics
async function loadDiagnostics() {
    const response = await apiRequest("/admin/status");
    if (response && response.ok) {
        const data = await response.json();
        updateStatusIndicator("db", data.database);
        updateStatusIndicator("redis", data.redis);
        updateStatusIndicator("wa", data.whatsapp_integration);
    }
}

function updateStatusIndicator(service, status) {
    const dot = document.getElementById(`status-${service}-dot`);
    const text = document.getElementById(`status-${service}-text`);

    dot.className = "dot";
    if (status === "connected" || status === "configured") {
        dot.classList.add("dot-green");
        text.innerText = status === "connected" ? "Online" : "Configurado";
    } else if (status === "disconnected") {
        dot.classList.add("dot-red");
        text.innerText = "Offline";
    } else {
        dot.classList.add("dot-red");
        text.innerText = "Não configurado";
    }
}

// List Agents
async function loadAgents() {
    const response = await apiRequest(`/admin/agents?limit=${agentsLimit}&offset=${agentsOffset}`);
    if (response && response.ok) {
        const agents = await response.json();
        const tbody = document.getElementById("agents-table-body");
        tbody.innerHTML = "";

        // Update page text
        const page = Math.floor(agentsOffset / agentsLimit) + 1;
        document.getElementById("page-indicator").innerText = `Página ${page}`;
        document.getElementById("btn-prev-page").disabled = agentsOffset === 0;
        document.getElementById("btn-next-page").disabled = agents.length < agentsLimit;

        agents.forEach(agent => {
            const tr = document.createElement("tr");

            const isSelf = agent.id === currentAdminId;
            const actionButtons = isSelf 
                ? `<span class="badge">Você (Atual)</span>`
                : `<div class="table-actions">
                     <button class="btn btn-secondary btn-sm" onclick="openEditModal('${agent.id}', '${agent.name}', '${agent.email}', '${agent.role}')">Editar</button>
                     <button class="btn btn-danger btn-sm" onclick="deleteAgent('${agent.id}')">Excluir</button>
                   </div>`;

            const lastActiveTime = agent.last_active ? new Date(agent.last_active).toLocaleString() : "Nunca";

            tr.innerHTML = `
                <td><strong>${agent.name}</strong></td>
                <td>${agent.email}</td>
                <td><span class="badge">${agent.role}</span></td>
                <td>${lastActiveTime}</td>
                <td class="text-right">${actionButtons}</td>
            `;
            tbody.appendChild(tr);
        });
    }
}

// Create Agent
async function handleCreateAgent(e) {
    e.preventDefault();
    const name = document.getElementById("create-name").value;
    const email = document.getElementById("create-email").value;
    const password = document.getElementById("create-password").value;
    const role = document.getElementById("create-role").value;

    const response = await apiRequest("/admin/agents", {
        method: "POST",
        body: JSON.stringify({ name, email, password, role })
    });

    if (response && response.ok) {
        showToast("Agente criado com sucesso!", "success");
        closeModal("modal-create");
        document.getElementById("create-agent-form").reset();
        loadAgents();
    } else if (response) {
        const data = await response.json().catch(() => ({}));
        showToast(data.error || "Falha ao criar agente.", "error");
    }
}

// Open Edit Modal
function openEditModal(id, name, email, role) {
    document.getElementById("edit-id").value = id;
    document.getElementById("edit-name").value = name;
    document.getElementById("edit-email").value = email;
    document.getElementById("edit-password").value = "";
    document.getElementById("edit-role").value = role;
    openModal("modal-edit");
}

// Update Agent
async function handleUpdateAgent(e) {
    e.preventDefault();
    const id = document.getElementById("edit-id").value;
    const name = document.getElementById("edit-name").value;
    const email = document.getElementById("edit-email").value;
    const password = document.getElementById("edit-password").value;
    const role = document.getElementById("edit-role").value;

    const payload = { name, email, role };
    if (password.trim() !== "") {
        payload.password = password;
    }

    const response = await apiRequest(`/admin/agents/${id}`, {
        method: "PUT",
        body: JSON.stringify(payload)
    });

    if (response && response.ok) {
        showToast("Agente atualizado com sucesso!", "success");
        closeModal("modal-edit");
        loadAgents();
    } else if (response) {
        const data = await response.json().catch(() => ({}));
        showToast(data.error || "Falha ao atualizar agente.", "error");
    }
}

// Delete Agent
async function deleteAgent(id) {
    if (!confirm("Tem certeza que deseja desativar este agente e revogar todas as suas sessões?")) {
        return;
    }

    const response = await apiRequest(`/admin/agents/${id}`, { method: "DELETE" });

    if (response && response.ok) {
        showToast("Agente desativado e sessões revogadas.", "success");
        loadAgents();
        loadSessions();
    } else if (response) {
        const data = await response.json().catch(() => ({}));
        showToast(data.error || "Falha ao deletar agente.", "error");
    }
}

// List Sessions
async function loadSessions() {
    const response = await apiRequest("/admin/sessions");
    if (response && response.ok) {
        const sessions = await response.json();
        const tbody = document.getElementById("sessions-table-body");
        tbody.innerHTML = "";

        sessions.forEach(sess => {
            const tr = document.createElement("tr");

            const isCurrentSession = sess.id === localStorage.getItem("session_id");
            const actionButton = isCurrentSession
                ? `<span class="badge">Sessão Atual</span>`
                : `<button class="btn btn-danger btn-sm" onclick="revokeSession('${sess.id}')">Revogar</button>`;

            const expiresTime = new Date(sess.expires_at).toLocaleString();

            tr.innerHTML = `
                <td><strong>${sess.agent_name || "Agente"}</strong></td>
                <td><code>${sess.ip_address || "Desconhecido"}</code></td>
                <td style="max-width: 200px; overflow: hidden; text-overflow: ellipsis;" title="${sess.user_agent || ''}">
                    ${sess.user_agent || 'N/A'}
                </td>
                <td>${expiresTime}</td>
                <td class="text-right">${actionButton}</td>
            `;
            tbody.appendChild(tr);
        });
    }
}

// Revoke Session
async function revokeSession(id) {
    if (!confirm("Tem certeza que deseja revogar esta sessão ativa e desconectar o agente?")) {
        return;
    }

    const response = await apiRequest(`/admin/sessions/${id}`, { method: "DELETE" });

    if (response && response.ok) {
        showToast("Sessão revogada com sucesso.", "success");
        loadSessions();
    } else if (response) {
        const data = await response.json().catch(() => ({}));
        showToast(data.error || "Falha ao revogar sessão.", "error");
    }
}

// ── Modal Utilities ──
function openModal(id) {
    document.getElementById(id).classList.remove("hidden");
}

function closeModal(id) {
    document.getElementById(id).classList.add("hidden");
}

// ── Toast Helper ──
function showToast(message, type = "success") {
    const container = document.getElementById("toast-container");
    const toast = document.createElement("div");
    toast.className = `toast toast-${type}`;
    toast.innerText = message;
    
    container.appendChild(toast);
    
    // Auto-remove toast after 4 seconds
    setTimeout(() => {
        toast.style.animation = "slideIn 0.3s reverse";
        setTimeout(() => {
            toast.remove();
        }, 300);
    }, 4000);
}

// Error logging helper
function slogError(message, err) {
    console.error(`[Admin UI] ${message}:`, err);
}
