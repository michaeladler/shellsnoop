{ lib, stdenv, buildGoModule, fetchFromGitHub, llvmPackages, libbpf, rev ? "" }:

buildGoModule rec {
  pname = "shellsnoop";
  version = "unstable-2023-10-12";

  src = ../.;

  nativeBuildInputs = [
    llvmPackages.llvm
    llvmPackages.clang
    libbpf
  ];

  buildInputs = [
    libbpf
  ];

  vendorHash = "sha256-cTtTCl1ku/D4alt1OrI7dWdzRuGfbBbLpbRJrIDRL7o=";

  ldflags = [
    "-s"
    "-w"
    "-X main.Commit=${rev}"
  ];

  CGO_ENABLED = 0;

  preBuild = ''
    export NIX_HARDENING_ENABLE=""
    go generate ./...
  '';

  postBuild = ''
    make shellsnoop-client
  '';

  postInstall = ''
    install -m0755 shellsnoop-client $out/bin/shellsnoop-client
  '';

  meta = with lib; {
    description = "eBPF program to snoop shell commands";
    license = licenses.asl20;
    maintainers = with maintainers; [ michaeladler ];
  };

}
