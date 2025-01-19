import { createRequire } from "module";
import { defineConfig, type DefaultTheme } from "vitepress";
import { tabsMarkdownPlugin } from "vitepress-plugin-tabs";

const require = createRequire(import.meta.url);
const pkg = require("../../package.json");

export default defineConfig({
  lang: "en-US",
  title: "Queryx",
  description: "Schema-first and type-safe ORM for Go and TypeScript",

  lastUpdated: true,
  cleanUrls: true,

  sitemap: {
    hostname: "https://queryx.caitouyun.com",
    transformItems(items) {
      return items.filter((item) => !item.url.includes("migration"));
    },
  },

  markdown: {
    config(md) {
      md.use(tabsMarkdownPlugin);
    },
  },

  head: [
    // ["link", { rel: "icon", href: "/vitepress-logo-mini.svg" }],
    ["meta", { name: "theme-color", content: "#5f67ee" }],
    ["meta", { name: "og:type", content: "website" }],
    ["meta", { name: "og:locale", content: "en" }],
    ["meta", { name: "og:site_name", content: "Queryx" }],
    // [
    //   "meta",
    //   { name: "og:image", content: "https://vitepress.dev/vitepress-og.jpg" },
    // ],
    // [
    //   "meta",
    //   {
    //     name: "twitter:image",
    //     content: "https://vitepress.dev/vitepress-og.jpg",
    //   },
    // ],
  ],

  themeConfig: {
    // logo: { src: "/vitepress-logo-mini.svg", width: 24, height: 24 },

    nav: nav(),

    sidebar: {
      "/docs/": { base: "/docs/", items: sidebarDocs() },
    },

    editLink: {
      pattern: "https://github.com/swiftcarrot/queryx/edit/main/website/:path",
      text: "Edit this page on GitHub",
    },

    socialLinks: [
      { icon: "github", link: "https://github.com/swiftcarrot/queryx" },
      { icon: "discord", link: "https://discord.gg/QUTxjJBRfA" },
    ],

    footer: {
      message: "Released under the Apache License.",
      copyright: "Copyright © 2021-present SWIFTCARROT Technologies",
    },

    search: {
      provider: "local",
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
      text: pkg.version,
      items: [
        {
          text: "Changelog",
          link: "https://github.com/swiftcarrot/queryx/blob/main/CHANGELOG.md",
        },
        {
          text: "Contributing",
          link: "https://github.com/swiftcarrot/queryx/blob/main/.github/contributing.md",
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
      ],
    },
    {
      text: "Tutorials",
      collapsed: false,
      items: [
        { text: "Query Methods", link: "query-methods" },
        { text: "Raw SQL", link: "raw-sql" },
        { text: "Transaction", link: "transaction" },
        { text: "Association", link: "association" },
        { text: "Data Types", link: "data-types" },
        { text: "Environment Variable", link: "environment-variable" },
        { text: "Database Index", link: "database-index" },
        { text: "Custom Table Name", link: "custom-table-name" },
        { text: "Custom Primary Key", link: "custom-primary-key" },
        { text: "Custom Database Time Zone", link: "time-zone" },
        { text: "Build from source", link: "build-from-source" },
        { text: "Docker Support", link: "docker-support" },
      ],
    },
    {
      text: "References",
      collapsed: false,
      items: [{ text: "CLI", link: "cli" }],
    },
  ];
}
