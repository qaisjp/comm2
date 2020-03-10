import { Injectable } from '@angular/core';

export interface User {
  readonly id: number;
  readonly created_at: Date;
  username: string;
}

export interface AuthenticatedUser extends User {
  readonly updated_at: Date;
}

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor() { }
}
