export type User = {
  id: number;
  name: string;
  created_at: string;
};

export type Employee = {
  id: number;
  name: string;
  email: string;
  created_at: string;
};

const defaultBase = "http://localhost:8080";

export function apiBase(): string {
  const raw = import.meta.env.VITE_API_BASE_URL ?? defaultBase;
  return raw.replace(/\/$/, "");
}

async function readBodyMessage(res: Response): Promise<string> {
  const text = (await res.text()).trim();
  return text || res.statusText;
}

export async function listUsers(): Promise<User[]> {
  const res = await fetch(`${apiBase()}/users`);
  if (!res.ok) {
    throw new Error(await readBodyMessage(res));
  }
  return res.json() as Promise<User[]>;
}

export async function createUser(name: string): Promise<User> {
  const res = await fetch(`${apiBase()}/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name }),
  });
  if (!res.ok) {
    throw new Error(await readBodyMessage(res));
  }
  return res.json() as Promise<User>;
}

export async function listEmployees(): Promise<Employee[]> {
  const res = await fetch(`${apiBase()}/employees`);
  if (!res.ok) {
    throw new Error(await readBodyMessage(res));
  }
  return res.json() as Promise<Employee[]>;
}

export async function createEmployee(
  name: string,
  email: string,
): Promise<Employee> {
  const res = await fetch(`${apiBase()}/employees`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email }),
  });
  if (!res.ok) {
    const msg = await readBodyMessage(res);
    const err = new Error(msg) as Error & { status?: number };
    err.status = res.status;
    throw err;
  }
  return res.json() as Promise<Employee>;
}
