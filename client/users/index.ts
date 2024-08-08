import { User, UserStatus } from "../../proto/users_pb";
import { allUsers } from "./services";

async function run() {
  const user = await getUser(1);
  // @ts-ignore
    console.log(user.toString());

  const jim = new User();
  jim.setName("Jim");
  jim.setAge(10);
  jim.setId(20);
  jim.setStatus(UserStatus.OFFLINE);
  createUsers([jim]);
  console.log(`\nCreated user ${jim.toString()}`);

  const users = await allUsers();
  console.log(`\nListing all ${users.length} users`);
  for (const user of users) {
    console.log(user.toString());
  }
}

run();

function getUser(arg0: number) {
    throw new Error("Function not implemented.");
}
function createUsers(arg0: User[]) {
    throw new Error("Function not implemented.");
}

