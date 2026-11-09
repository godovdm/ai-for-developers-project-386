// Apply compatibility fixes for exactOptionalPropertyTypes to generated SDK code.
import { readFileSync, writeFileSync } from 'node:fs';

const path = process.argv[2];
if (!path) throw new Error('Expected path to generated runtime.ts');
let code = readFileSync(path, 'utf8');
const replacements = [
  ['this.configuration = configuration;', 'this.configuration = configuration.configuration;'],
  ['preMiddlewares.map((pre) => ({ pre }))', 'preMiddlewares.filter((pre): pre is NonNullable<typeof pre> => pre !== undefined).map((pre) => ({ pre }))'],
  ['postMiddlewares.map((post) => ({ post }))', 'postMiddlewares.filter((post): post is NonNullable<typeof post> => post !== undefined).map((post) => ({ post }))'],
  ['credentials: this.configuration.credentials,', '...(this.configuration.credentials === undefined ? {} : { credentials: this.configuration.credentials }),'],
  ['response?: Response;', 'response?: Response | undefined;'],
];
for (const [from, to] of replacements) {
  if (!code.includes(from)) throw new Error(`Generator output changed: ${from}`);
  code = code.replace(from, to);
}
writeFileSync(path, code);
