import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

@Injectable()
export class TodoService {
  constructor(private httpClient: HttpClient) {}

  getTodoList() {
    return this.httpClient.get(environment.gateway + '/api/todo');
  }

  addTodo(todo: Todo) {
    return this.httpClient.post(environment.gateway + '/api/todo', todo);
  }

  completeTodo(todo: Todo) {
    return this.httpClient.put(environment.gateway + '/api/todo', todo);
  }

  deleteTodo(todo: Todo) {
    return this.httpClient.delete(environment.gateway + '/api/todo/' + todo.id);
  }
}

export class Todo {
  id: number;
  message: string;
  complete: boolean;
}
