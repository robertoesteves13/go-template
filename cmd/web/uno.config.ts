import { defineConfig, presetAttributify, presetWind3 } from 'unocss'
import fs from 'fs';

const content = fs.readFileSync('./node_modules/@unocss/reset/tailwind.css').toString();

export default defineConfig({
  presets: [
    presetWind3(),
    presetAttributify(),
  ],
  preflights: [
    {
      getCSS: () => content,
    },
  ],
  cli: {
    entry: {
      patterns: ['templates/**/*.templ'],
      outFile: 'global.css'
    }
  },
})

