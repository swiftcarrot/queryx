// import { createRequire } from "module";
import { defineConfig, type DefaultTheme } from "vitepress";

// const require = createRequire(import.meta.url);
// const pkg = require("vitepress/package.json");

export default defineConfig({
  lang: "en-US",
  title: "Queryx",
  description: "Schema-first and type-safe ORM for Go and TypeScript",

  lastUpdated: true,
  cleanUrls: true,

  sitemap: {
    hostname: "https://vitepress.dev",
    transformItems(items) {
      return items.filter((item) => !item.url.includes("migration"));
    },
  },

  head: [
    ["link", { rel: "icon", href: "/vitepress-logo-mini.svg" }],
    ["meta", { name: "theme-color", content: "#5f67ee" }],
    ["meta", { name: "og:type", content: "website" }],
    ["meta", { name: "og:locale", content: "en" }],
    ["meta", { name: "og:site_name", content: "VitePress" }],
    [
      "meta",
      { name: "og:image", content: "https://vitepress.dev/vitepress-og.jpg" },
    ],
    [
      "meta",
      {
        name: "twitter:image",
        content: "https://vitepress.dev/vitepress-og.jpg",
      },
    ],
    [
      "script",
      {
        src: "https://cdn.usefathom.com/script.js",
        "data-site": "AZBRSFGG",
        "data-spa": "auto",
        defer: "",
      },
    ],
  ],

  themeConfig: {
    // logo: { src: "/vitepress-logo-mini.svg", width: 24, height: 24 },

    nav: nav(),

    sidebar: {
      "/docs/": { base: "/docs/", items: sidebarDocs() },
    },

    editLink: {
      pattern: "https://github.com/swiftcarrot/queryx/edit/main/docs/:path",
      text: "Edit this page on GitHub",
    },

    socialLinks: [
      { icon: "github", link: "https://github.com/swiftcarrot/queryx" },
    ],

    footer: {
      message: "Released under the Apache License.",
      copyright: "Copyright © 2023-present SWIFTCARROT Technologies",
    },

    // search: {
    //   provider: "algolia",
    //   options: {
    //     appId: "8J64VVRP8K",
    //     apiKey: "a18e2f4cc5665f6602c5631fd868adfd",
    //     indexName: "vitepress",
    //   },
    // },
  },
});

function nav(): DefaultTheme.NavItem[] {
  return [
    {
      text: "Documentation",
      link: "/docs/what-is-queryx",
      activeMatch: "/docs/",
    },
    {
      // text: pkg.version,
      items: [
        {
          text: "Changelog",
          link: "https://github.com/vuejs/vitepress/blob/main/CHANGELOG.md",
        },
        {
          text: "Contributing",
          link: "https://github.com/vuejs/vitepress/blob/main/.github/contributing.md",
        },
      ],
    },
  ];
}

/* prettier-ignore */
function sidebarDocs(): DefaultTheme.SidebarItem[] {
  return [
    {
      text: "Introduction",
      collapsed: false,
      items: [
        { text: "What is Queryx?", link: "what-is-queryx" },
        { text: "Getting Started", link: "getting-started" },
        {
          text: "Database Connection",
          link: "database-connection",
          collapsed: true,
          items: [
            { text: "Go", link: "go" },
            { text: "TypeScript", link: "typescript" },
          ],
        },
      ],
    },
    {
      text: "Schema",
      collapsed: false,
      items: [
        { text: "Data Types", link: "data-types" },
        { text: "Primary Key", link: "primary-key" },
        { text: "Database Index", link: "database-index" },
        { text: "Association", link: "association" },
      ],
    },
    {
      text: "Database Management",
      collapsed: false,
      items: [
        { text: "Migrate", link: "migrate" }
      ],
    },
    {
      text: "ORM Methods",
      collapsed: false,
      items: [
        { text: "Transaction", link: "transaction" }
      ],
    },
    {
      text: "References",
      collapsed: false,
      items: [
        { text: 'Data Types', link: "data-types" },
        { text: 'CLI', link: "CLI" }
      ]
    },
  ];
}
