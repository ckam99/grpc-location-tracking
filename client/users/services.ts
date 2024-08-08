import { Empty } from "google-protobuf/google/protobuf/empty_pb";
import { User, UserRequest } from "../../proto/users_pb";
import { client, noop } from "./utils";

export  function allUsers() {
  return new Promise<User[]>((resolve, reject) => {
    const stream = client.getUsers(new Empty());
    const users: User[] = [];
    stream.on("data", (user) => users.push(user));
    stream.on("error", reject);
    stream.on("end", () => resolve(users));
  });
}

export  function createNewUsers(users: User[]) {
  const stream = client.createUser(noop);
  for (const user of users) {
    stream.write(user);
  }
  stream.end();
}

export  function getUsers(id: number) {
  return new Promise<User>((resolve, reject) => {
    const request = new UserRequest();
    request.setId(id);

    client.getUser(request, (err, user) => {
      if (err) reject(err);
      else resolve(user);
    });
  });
}