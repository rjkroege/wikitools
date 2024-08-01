export type { WikiArticleParms };

type WikiArticleParms = {
  title: string;
  template: string;
};

export function ShowNewArticleDialog(
  templates: string[],
): Promise<WikiArticleParms> {
  const thealert = new Alert();
  thealert.title = "New Wiki Article";
  const tf = thealert.addTextField("Article title", "thoughts");
  // I can use tf to alter the font or the like. I'm not sure that I need to.

  templates.forEach((e) => thealert.addAction(e));

  // Show the alert and wait for its dismisal
  return thealert.presentAlert().then((val: number): WikiArticleParms => {
    return {
      title: tf.text,
      template: templates[val],
    };
  });
}
