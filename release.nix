{ lib, stdenv, buildGoModule, fetchFromGitHub, llvmPackages, libbpf, rev ? "" }:

buildGoModule rec {
  pname = "shellsnoop";
  version = "0.1.0";

  src = ./.;

  nativeBuildInputs = [
    llvmPackages.llvm
    llvmPackages.clang
    libbpf
  ];

  buildInputs = [
    libbpf
  ];

  vendorHash = "sha256-ZjCkLcVm8JWXpJYvIYUr+f4LZyprmwfulePNpvz0vOk=";

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
