import React, { useEffect, useState, useRef } from 'react';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

const BASE_GITHUB_URL = 'https://github.com/mify-io/mify/releases'

async function fetchLatestReleaseAssetNames(owner, repo) {
  const url = `https://api.github.com/repos/${owner}/${repo}/releases/latest`;

  const response = await fetch(url, {
    headers: {
      'Accept': 'application/vnd.github.v3+json',
    }
  });

  if (!response.ok) {
    throw new Error(`GitHub API responded with status ${response.status}`);
  }

  const release = await response.json();
  const version = release.tag_name || release.name;
  return version.startsWith('v') ? version.slice(1) : version;
}

function pkglink(target, version, arch, full = true) {
    const archDict = {
        "arch": {
            "amd64": "x86_64",
            "arm64": "aarch64",
        },
        "rpm": {
            "amd64": "x86_64",
            "arm64": "aarch64",
        },
        "linux": {
            "amd64": "x86_64",
        },
        "windows": {
            "amd64": "x86_64",
        },
        "deb": {},
        "mac": {},
    }
    const archStr = archDict[target][arch] ?? arch;
    const targetsDict = {
        "arch": `mify-${version}-1-${archStr}.pkg.tar.zst`,
        "rpm": `mify-${version}-1.${archStr}.rpm`,
        "mac": `mify-darwin-all.tar.gz`,
        "linux": `mify-linux-${archStr}.tar.gz`,
        "windows": `mify-windows-${archStr}.zip`,
        "deb": `mify_${version}_${archStr}.deb`,
    }
    if (full) {
        return BASE_GITHUB_URL + "/download/v" + version + "/" + targetsDict[target]
    }
    return targetsDict[target]
}

function BinaryDownloadBlock({target, version, arch}) {
    return (
        <div class="download-block">
        <div class="download-info">
        <b>{arch.toUpperCase()}</b><br/>
        <b>Version:</b> {version}
        </div>
        <div class="download-button"><a href={pkglink(target, version, arch)}>Download</a></div>
        </div>
    )
}

export default function InstallGuide() {
  const [version, setVersion] = useState('');
  const componentMounted = useRef(true);

  useEffect(() => {
    (async () => {
      var tag = await fetchLatestReleaseAssetNames('mify-io', 'mify')
      if (componentMounted.current) {
          setVersion(tag);
      }
      return () => {
          componentMounted.current = false;
      }
    })();
  }, []);  // The empty array ensures this useEffect runs once, similar to componentDidMount
  return (
  <div class="installation">
  <Tabs groupId="operating-systems">
  <TabItem value="mac" label="macOS">
  <CodeBlock>{
`brew tap mify-io/tap
brew install mify-io/tap/mify
`}
  </CodeBlock>

  <b>Binary downloads</b>
  <BinaryDownloadBlock target="mac" version={version} arch="amd64"/>
  <BinaryDownloadBlock target="mac" version={version} arch="arm64"/>

  </TabItem>
  <TabItem value="linux" label="Linux">
  <Tabs groupId="linux-distros">
  <TabItem value="deb" label="Ubuntu/Debian">
  <b>AMD64</b>
  <CodeBlock>{
`wget ${pkglink("deb", version, "amd64")}
sudo dpkg -i ${pkglink("deb", version, "amd64", false)}
`}
  </CodeBlock>

  <b>ARM64</b>
  <CodeBlock>{
`wget ${pkglink("deb", version, "arm64")}
sudo dpkg -i ${pkglink("deb", version, "arm64", false)}
`}
  </CodeBlock>

  </TabItem>
  <TabItem value="rpm" label="CentOS/RHEL/Fedora">
  <b>AMD64</b>
  <CodeBlock>{
`wget ${pkglink("rpm", version, "amd64")}
sudo rpm -i ${pkglink("rpm", version, "amd64", false)}
`}
  </CodeBlock>

  <b>ARM64</b>
  <CodeBlock>{
`wget ${pkglink("rpm", version, "arm64")}
sudo rpm -i ${pkglink("rpm", version, "arm64", false)}
`}
  </CodeBlock>

  </TabItem>
  <TabItem value="arch" label="Arch">
  <b>AMD64</b>
  <CodeBlock>{
`wget ${pkglink("arch", version, "amd64")}
sudo pacman -U ${pkglink("arch", version, "amd64", false)}
`}
  </CodeBlock>

  <b>ARM64</b>
  <CodeBlock>{
`wget ${pkglink("arch", version, "arm64")}
sudo pacman -U ${pkglink("arch", version, "arm64", false)}
`}
  </CodeBlock>
  </TabItem>
  <TabItem value="brew" label="Homebrew">
  <CodeBlock>{
`brew tap mify-io/tap
brew install mify-io/tap/mify
`}
  </CodeBlock>
  </TabItem>
  </Tabs>

  <b>Binary downloads</b>
  <BinaryDownloadBlock target="linux" version={version} arch="amd64"/>
  <BinaryDownloadBlock target="linux" version={version} arch="arm64"/>
  </TabItem>
  <TabItem value="win" label="Windows">
  <b>Binary downloads</b>
  <BinaryDownloadBlock target="windows" version={version} arch="amd64"/>
  <BinaryDownloadBlock target="windows" version={version} arch="arm64"/>
  </TabItem>
</Tabs>
  <b>Release information</b>
  <div class="download-block">
  <div class="download-info">
  <b>Changelog</b><br/>
  <b>Version:</b> {version}
  </div>
  <div class="download-button"><a href={BASE_GITHUB_URL + "/tag/v" + version}>GitHub</a></div>
  </div>

</div>
);

}
