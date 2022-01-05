import { readFileSync, writeFile } from "fs";
import {
  resolveSystem,
  listSystems,
  JSON_PATH,
  getSystemByPrefix,
  listPrefixes,
} from "./config.js";
const readJson = (path) => JSON.parse(readFileSync(path, "utf-8"));

const obj = readJson("./fc2_all_games.json");
const systems = listSystems();

const stripPrefix = (game, prefix) =>
  game.gameid.split(`${prefix}_`)[1] || game.gameid;
let totalMissing = 0;

const getGameByPrefix = (prefix) => {
  const sys = getSystemByPrefix(prefix);
  if (sys.emulator !== "fbneo") {
    return obj.filter((game) => game.emulator === sys.emulator);
  }
  if (sys.emulator === "fbneo" && prefix === "") {
    return obj.filter(
      (game) => game.emulator === "fbneo" && game.gameid.split("_").length <= 1
    );
  }
  if (sys.emulator === "fbneo" && prefix !== "") {
    return obj.filter(
      (game) =>
        game.emulator === "fbneo" && game.gameid.startsWith(`${prefix}_`)
    );
  }
};

const writeToJson = (source, filename) => {
  writeFile(
    `./output/${filename}.json`,
    JSON.stringify(source, null, 4),
    { encoding: "utf-8" },
    () => null
  );
};

const getGames = (prefix, filepath) => {
  const system = getSystemByPrefix(prefix);
  const games = Object.keys(readJson(filepath || JSON_PATH + system.jsonFile));

  return games;
};

const getMissingGames = (prefix) => {
  const system = getSystemByPrefix(prefix);
  const allGames = getGameByPrefix(prefix);
  const games = getGames(prefix);
  const missingGames = allGames.filter((item) => {
    switch (system.prefix) {
      case "snes":
      case "nes":
      case "flycast":
        return !games.includes(item.gameid);
      default:
        return !games.includes(stripPrefix(item, prefix));
    }
  });

  console.log(`${system.jsonFile} (${system.emulator})
    Missing games: ${missingGames.length}
    `);

  totalMissing += missingGames.length;

  writeToJson(missingGames, system.jsonFile);
};

// getMissingGames("flycast");
console.log(listPrefixes());

const main = () => {
  const prefixes = listPrefixes();

  prefixes.map((prefix) => getMissingGames(prefix));
  console.log(`Done: missing ${totalMissing} games`);
};

main();
