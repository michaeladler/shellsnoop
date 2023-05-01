inputs: { config, lib, utils, pkgs, ... }:

with lib;

let
  cfg = config.programs.shellsnoop;
  shellsnoopPkg = inputs.self.packages.${pkgs.stdenv.hostPlatform.system}.default;
in
{
  options.programs.shellsnoop = {
    enable =
      mkEnableOption null
      // {
        description = "Enable shellsnoop";
      };

    package = mkOption {
      type = types.path;
      default = shellsnoopPkg;
      defaultText = literalExpression ''
        shellsnoopPkg
      '';
      example = literalExpression "pkgs.shellsnoop";
      description = mdDoc ''
        The shellsnoop package to use.
      '';
    };

    uid = mkOption {
      type = types.int;
      default = 1000;
      description = lib.mdDoc ''
        The UID for which to snoop shell commands.
      '';
    };

    socket = mkOption {
      type = types.str;
      default = "/run/shellsnoop/shellsnoop.sock";
      description = lib.mdDoc ''
        Path for the unix socket file (used by shellsnoop-client)
      '';
      example = "${runtimeDir}/<name>.sock";
    };

  };

  config = mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];

    systemd.services.shellsnoop = {
      description = "snoop shell commands";
      wantedBy = [ "multi-user.target" ];
      partOf = [ "multi-user.target" ];
      serviceConfig = {
        ExecStart = "${cfg.package}/bin/shellsnoop -u ${builtins.toString cfg.uid} -s ${utils.escapeSystemdExecArg cfg.socket}";
        RuntimeDirectory = "shellsnoop";
      };
    };
  };
}
