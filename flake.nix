{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = inputs @ { self, nixpkgs }:

    let
      shortRev = if (self ? shortRev) then self.shortRev else "dev";

      systems = [
        "x86_64-linux"
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);

      mkPkgs = system: import nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      };

    in

    {
      packages = forAllSystems (system:
        let pkgs = mkPkgs system; in
        {

          inherit (pkgs) shellsnoop;
          default = pkgs.shellsnoop;
        });

      devShells = forAllSystems
        (system:
          let pkgs = mkPkgs system; in
          {
            default = with pkgs; mkShell {
              nativeBuildInputs = [ go ];
              buildInputs = [ libbpf ];
              hardeningDisable = [ "all" ];
            };
          });

      nixosModules.default = import ./nix/module.nix inputs;


      overlays.default = _: prev: {
        shellsnoop = prev.callPackage ./nix/release.nix {
          rev = shortRev;
          llvmPackages = prev.llvmPackages_16;
        };
      };
    };
}
