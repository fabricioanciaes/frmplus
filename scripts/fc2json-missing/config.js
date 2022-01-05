export const JSON_PATH = "./jsonFiles/";
export const EMULATOR_INITIALS = {
  "Arcade FC2": {
    emulator: "fbneo",
    prefix: "",
    jsonFile: "fbneo_roms.json",
  },
  NAOMI: {
    emulator: "flycast",
    prefix: "flycast",
    jsonFile: "flycast_roms.json",
  },
  "Arcade FC1": {
    emulator: "fc1",
    prefix: "fc1",
    jsonFile: "fc1_roms.json",
  },
  Dreamcast: {
    emulator: "flycast",
    prefix: "flycast",
    jsonFile: "flycast_roms.json",
  },
  Atomiswave: {
    emulator: "flycast",
    prefix: "flycast",
    jsonFile: "flycast_roms.json",
  },
  "Super NES": {
    emulator: "snes9x",
    prefix: "snes",
    jsonFile: "snes9x_roms.json",
  },
  Megadrive: {
    emulator: "fbneo",
    prefix: "md",
    jsonFile: "fbneo_md_roms.json",
  },
  NES: {
    emulator: "fbneo",
    prefix: "nes",
    jsonFile: "fbneo_nes_roms.json",
  },
  "MSX 1": {
    emulator: "fbneo",
    prefix: "msx",
    jsonFile: "fbneo_msx_roms.json",
  },
  TurboGrafx16: {
    emulator: "fbneo",
    prefix: "tg",
    jsonFile: "fbneo_tg_roms.json",
  },
  "Master System": {
    emulator: "fbneo",
    prefix: "sms",
    jsonFile: "fbneo_sms_roms.json",
  },
  "Sega SG-1000": {
    emulator: "fbneo",
    prefix: "sg1k",
    jsonFile: "fbneo_sg1k_roms.json",
  },
  "PC-Engine": {
    emulator: "fbneo",
    prefix: "pce",
    jsonFile: "fbneo_pce_roms.json",
  },
  ColecoVision: {
    emulator: "fbneo",
    prefix: "cv",
    jsonFile: "fbneo_cv_roms.json",
  },
  "Game Gear": {
    emulator: "fbneo",
    prefix: "gg",
    jsonFile: "fbneo_gg_roms.json",
  },
};

export const resolveSystem = (system) => EMULATOR_INITIALS[system];
export const getSystemByPrefix = (prefix) =>
  listSystems().find((item) => item.prefix === prefix);
export const listSystems = () => {
  return Object.keys(EMULATOR_INITIALS).map((item) => resolveSystem(item));
};

export const listPrefixes = () => {
  const list = new Set();
  Object.keys(EMULATOR_INITIALS).map((item) =>
    list.add(resolveSystem(item).prefix)
  );

  return Array.from(list);
};
