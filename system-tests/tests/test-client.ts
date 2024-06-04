export interface User {
  id: string;
  name: string;
}

export class Client {
  private url: string;

  constructor(url: string) {
    this.url = url;
  }

  async getUsers(): Promise<User[]> {
    const res = await fetch(`${this.url}/`);
    const { users } = await res.json();
    return users;
  }

  async deleteUser(id: string): Promise<void> {
    await fetch(`${this.url}/id/${id}`, {
      method: 'DELETE',
    });
  }
}