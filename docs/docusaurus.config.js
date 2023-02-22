// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Mify',
  tagline: 'Microservice infrastructure for you',
  url: 'https://mify.io',
  baseUrl: '/docs/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'facebook', // Usually your GitHub org/user name.
  projectName: 'docusaurus', // Usually your repo name.

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/mify-io/mify/tree/main/docs/',
          routeBasePath: '/',
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        // title: 'Documentation',
        logo: {
          alt: 'Mify Logo',
          src: 'img/logo.png',
          srcDark: 'img/logo-white.png',
          href: 'https://mify.io',
          target: '_self',
        },
        items: [
          {
            type: 'doc',
            docId: 'getting-started/index',
            position: 'right',
            label: 'Docs',
          },
          {
            type: 'doc',
            docId: 'cloud/overview',
            position: 'right',
            label: 'Cloud',
          },
          {
            href: 'https://mify.io/pricing',
            label: 'Pricing',
            position: 'right',
          },
          {
            href: 'https://github.com/mify-io/mify',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Get started',
                to: '/docs/',
              },
              {
                label: 'Create Service',
                to: '/docs/guides/overview',
              },
              {
                label: 'Deploy to Cloud',
                to: '/docs/cloud/overview',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Slack',
                href: 'https://join.slack.com/t/mifyio/shared_invite/zt-1llnbiio6-lG45E696QOEVzHb0__Qqxw',
              },
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/mify',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/mify-io/mify',
              },
            ],
          },
        ],
        copyright: `© Copyright ${new Date().getFullYear()} Mify`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
