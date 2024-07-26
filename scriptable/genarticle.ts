// Generate a complete wikitools article

export function GenerateArticle(
  template: string,
  title: string,
  date: string,
): string {
  const output: string = `---
title: ${title}
date: ${date}
tags: 
---

${template}
`;

  return output;
}
