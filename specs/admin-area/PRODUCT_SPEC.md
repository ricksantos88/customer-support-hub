# Product Specification - Administrative Area

## Feature Objectives
Provide system administrators with a secure console/API to manage agents, active sessions, and oversee the customer support system. This area allows managing the lifecycle of support staff and revoking compromised sessions.

## Functional Requirements
1. **Agent Management (CRUD)**:
   - Create new agents with designated roles (`admin` or `agent`).
   - List all agents currently registered in the system.
   - Update agent details (Name, Email, Password, Role).
   - Deactivate/Remove agents from the system to deny access.
2. **Session Monitoring and Revocation**:
   - List all active agent sessions (revealing Agent Name, IP Address, User Agent, Last Active Time, Expiry Time).
   - Revoke/Terminate an active session remotely (forces the target agent to log out).
3. **Admin Dashboard Metrics (Metadata/Status)**:
   - Provide basic system diagnostics (DB connectivity status, Redis status, WhatsApp config setup verification).

## Business Rules
* **Strict Admin Access Only**: Only agents authenticated with the `admin` role are permitted to view or execute operations under the administrative endpoints.
* **Non-Self-Demotion**: An admin cannot demote themselves to the `agent` role.
* **Non-Self-Deletion**: An admin cannot delete/deactivate their own account to prevent lockout.
* **Immediate Revocation Effect**: When a session is revoked, the cached session in Redis must be invalidated instantly, and subsequent requests using that session's access token must return HTTP 401 Unauthorized.

## Acceptance Criteria
* Authenticating with an `agent` role credentials and calling any `/admin` endpoint returns HTTP 403 Forbidden.
* Authenticating with an `admin` role credentials allows successful execution of admin CRUD endpoints.
* Deleting/deactivating an agent immediately revokes all of their active sessions.
* Revoking a session by ID successfully logs out the target agent session.
* Non-self-deletion and non-self-demotion rules must return HTTP 400 Bad Request with a clear message.
