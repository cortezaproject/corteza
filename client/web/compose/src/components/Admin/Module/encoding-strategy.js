export const types = {
  Omit: 'omit',
  Plain: 'plain',
  Alias: 'alias',
  JSON: 'json',
}

// default for system fields means null!
export function systemFieldStrategyConfig (strategy, config) {
  switch (strategy) {
    case types.JSON:
      return { [strategy]: { ident: config.ident } }
    case types.Alias:
      return { [strategy]: { ident: config.ident } }
    case types.Omit:
      return { [strategy]: true }

    default:
      return null
  }
}

// default for module fields means JSON!
export function moduleFieldStrategyConfig (strategy, config) {
  switch (strategy) {
    case types.Plain:
      return { [types.Plain]: {} }
    case types.Alias:
      return { [strategy]: { ident: config.ident } }
    case types.JSON:
      return { [strategy]: { ident: config.ident } }
    case types.Omit:
      return { [strategy]: true }

    default:
      return null
  }
}

// extract ident from config
export function defaultConfigDraft (config, ident) {
  if (typeof config !== 'object') {
    return { ident, def: 1 }
  }

  for (const t of [types.Alias, types.JSON]) {
    if (config[t] && config[t].ident) {
      return { ...config[t] }
    }
  }

  return { ident, def: 2 }
}
