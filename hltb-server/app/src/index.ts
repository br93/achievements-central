import { Elysia } from "elysia";
import { hltb } from "./hltb";

const app = new Elysia()
    .get("/hltb/:game", async ({params}) => {
      return await hltb.search(params.game)
    })
   
  .listen(3000);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
