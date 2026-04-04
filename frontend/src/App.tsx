import { useCallback, useEffect, useState } from "react";
import {
  createEmployee,
  createUser,
  listEmployees,
  listUsers,
  type Employee,
  type User,
} from "./api";

type Section = "users" | "employees";

function formatTime(iso: string): string {
  try {
    return new Date(iso).toLocaleString();
  } catch {
    return iso;
  }
}

export default function App() {
  const [section, setSection] = useState<Section>("users");

  const [users, setUsers] = useState<User[]>([]);
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [usersLoading, setUsersLoading] = useState(true);
  const [employeesLoading, setEmployeesLoading] = useState(true);

  const [userName, setUserName] = useState("");
  const [empName, setEmpName] = useState("");
  const [empEmail, setEmpEmail] = useState("");

  const [usersError, setUsersError] = useState<string | null>(null);
  const [employeesError, setEmployeesError] = useState<string | null>(null);
  const [userFormError, setUserFormError] = useState<string | null>(null);
  const [empFormError, setEmpFormError] = useState<string | null>(null);

  const loadUsers = useCallback(async () => {
    setUsersLoading(true);
    setUsersError(null);
    try {
      setUsers(await listUsers());
    } catch (e) {
      setUsersError(e instanceof Error ? e.message : "Failed to load users");
    } finally {
      setUsersLoading(false);
    }
  }, []);

  const loadEmployees = useCallback(async () => {
    setEmployeesLoading(true);
    setEmployeesError(null);
    try {
      setEmployees(await listEmployees());
    } catch (e) {
      setEmployeesError(
        e instanceof Error ? e.message : "Failed to load employees",
      );
    } finally {
      setEmployeesLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadUsers();
    void loadEmployees();
  }, [loadUsers, loadEmployees]);

  async function onCreateUser(e: React.FormEvent) {
    e.preventDefault();
    setUserFormError(null);
    try {
      await createUser(userName.trim());
      setUserName("");
      await loadUsers();
    } catch (err) {
      setUserFormError(
        err instanceof Error ? err.message : "Failed to create user",
      );
    }
  }

  async function onCreateEmployee(e: React.FormEvent) {
    e.preventDefault();
    setEmpFormError(null);
    try {
      await createEmployee(empName.trim(), empEmail.trim());
      setEmpName("");
      setEmpEmail("");
      await loadEmployees();
    } catch (err) {
      const msg =
        err instanceof Error ? err.message : "Failed to create employee";
      setEmpFormError(msg);
    }
  }

  return (
    <div className="app">
      <header className="header">
        <h1>User &amp; employee management</h1>
        <p className="muted">
          API: <code>{import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080"}</code>
        </p>
        <nav className="tabs" aria-label="Section">
          <button
            type="button"
            className={section === "users" ? "tab active" : "tab"}
            onClick={() => setSection("users")}
          >
            Users
          </button>
          <button
            type="button"
            className={section === "employees" ? "tab active" : "tab"}
            onClick={() => setSection("employees")}
          >
            Employees
          </button>
        </nav>
      </header>

      {section === "users" && (
        <section className="panel" aria-labelledby="users-heading">
          <h2 id="users-heading">Users</h2>

          <form className="form" onSubmit={onCreateUser}>
            <label htmlFor="user-name">Name</label>
            <div className="form-row">
              <input
                id="user-name"
                name="name"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                placeholder="Jane Doe"
                autoComplete="name"
                required
              />
              <button type="submit">Create user</button>
            </div>
            {userFormError && (
              <p className="error" role="alert">
                {userFormError}
              </p>
            )}
          </form>

          {usersError && (
            <p className="error" role="alert">
              {usersError}
            </p>
          )}
          {usersLoading ? (
            <p className="muted">Loading users…</p>
          ) : (
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Created</th>
                  </tr>
                </thead>
                <tbody>
                  {users.length === 0 ? (
                    <tr>
                      <td colSpan={3} className="muted">
                        No users yet.
                      </td>
                    </tr>
                  ) : (
                    users.map((u) => (
                      <tr key={u.id}>
                        <td>{u.id}</td>
                        <td>{u.name}</td>
                        <td>{formatTime(u.created_at)}</td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          )}
        </section>
      )}

      {section === "employees" && (
        <section className="panel" aria-labelledby="employees-heading">
          <h2 id="employees-heading">Employees</h2>

          <form className="form" onSubmit={onCreateEmployee}>
            <label htmlFor="emp-name">Name</label>
            <input
              id="emp-name"
              name="empName"
              value={empName}
              onChange={(e) => setEmpName(e.target.value)}
              placeholder="Bob"
              autoComplete="name"
              required
            />
            <label htmlFor="emp-email">Email</label>
            <div className="form-row">
              <input
                id="emp-email"
                name="email"
                type="email"
                value={empEmail}
                onChange={(e) => setEmpEmail(e.target.value)}
                placeholder="bob@example.com"
                autoComplete="email"
                required
              />
              <button type="submit">Create employee</button>
            </div>
            {empFormError && (
              <p
                className={
                  empFormError.includes("email already exists")
                    ? "warn"
                    : "error"
                }
                role="alert"
              >
                {empFormError}
              </p>
            )}
          </form>

          {employeesError && (
            <p className="error" role="alert">
              {employeesError}
            </p>
          )}
          {employeesLoading ? (
            <p className="muted">Loading employees…</p>
          ) : (
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Created</th>
                  </tr>
                </thead>
                <tbody>
                  {employees.length === 0 ? (
                    <tr>
                      <td colSpan={4} className="muted">
                        No employees yet.
                      </td>
                    </tr>
                  ) : (
                    employees.map((e) => (
                      <tr key={e.id}>
                        <td>{e.id}</td>
                        <td>{e.name}</td>
                        <td>{e.email}</td>
                        <td>{formatTime(e.created_at)}</td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          )}
        </section>
      )}
    </div>
  );
}
