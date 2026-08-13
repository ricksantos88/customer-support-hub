# Product Specification - Administrative Web Interface

## Feature Objectives
Provide a premium, modern, and user-friendly Web Admin Dashboard served directly by the backend to manage agents, monitor active sessions, and check overall system diagnostics without needing command-line tools or raw SQL access.

## Functional Requirements
1. **Login View**:
   - Secure login form for administrators using email and password.
   - Redirect to dashboard upon successful login.
   - Local token caching (secure storage of access and refresh tokens).
2. **Dashboard Overview Panel**:
   - Diagnostic status metrics (Database: connected/disconnected, Redis: connected/disconnected, WhatsApp: configured/unconfigured).
   - Summary counts (Active Sessions, Total Agents).
3. **Agent Management Panel (CRUD)**:
   - Paginated list of registered agents showing Name, Email, Role, and Last Active status.
   - Create Agent Modal: Interactive form to register a new agent (name, email, password, role).
   - Edit Agent Modal: Interactive form to update agent profile details.
   - Delete/Inactivate Agent: Single-click action with confirmation prompt.
4. **Session Management Panel**:
   - Table of active sessions showing Agent Name, IP, User Agent, and login duration.
   - Revoke Session: Immediate termination of a session via click.
5. **Auto-logout & Security Checks**:
   - Frontend validation checking token expiry, with redirect to login on expiry.

## Business Rules
- Only users with the `admin` role can access the dashboard.
- Attempts to access the dashboard or make API requests without valid admin credentials automatically clear local storage and redirect to the login screen.
- Self-deactivation and self-demotion are blocked at both UI and API levels.

## Acceptance Criteria
- Loading `http://localhost:8080/admin-panel/` with no credentials shows the Login view.
- Submitting the login form with valid admin credentials successfully transitions to the Dashboard view.
- Submitting the login form with agent-level credentials displays a clear "Access Forbidden - Administrators Only" error message.
- Clicking "Create Agent" adds the agent to the table and displays a success toast.
- Clicking "Revoke Session" terminates the session in the database and updates the sessions list immediately.
