use crate::error::{Error, ErrorKind, Result};
use clap::{arg, ArgAction, ArgGroup, ArgMatches, Command};
use serde::{Deserialize, Serialize};
use std::fs;

/// Config contains all infos which components use.
#[derive(Serialize, Deserialize, Debug, Clone)]
#[serde(rename_all = "kebab-case")]
pub struct Config {
    pub auth_info: AuthInfo,
    pub server_info: ServerInfo,
}

impl Config {
    /// Load to build `Config` from file with yaml.
    ///
    /// # Example
    ///
    /// ```no_run
    /// let config = Config::load_from_file("./config_file.yaml").unwrap();
    /// println!("{:?}", config);
    /// ```
    pub fn load_from_file(path: &str) -> Result<Self> {
        let s = fs::read_to_string(path)
            .map_err(|e| Error::new_with_str(ErrorKind::OpenFile, e.to_string()))?;
        let r: Self = serde_yaml::from_str(s.as_str())
            .map_err(|e| Error::new_with_str(ErrorKind::ParseError, e.to_string()))?;
        Ok(r)
    }

    /// Add flags to main App's Command/Subcommand.
    ///
    /// Note that the ownership of the cmd move into and return in the end of this function.
    pub fn add_to_command(cmd: Command) -> Command {
        const ID: &str = "iam-config";
        cmd.group(ArgGroup::new(ID).multiple(true))
            .next_help_heading("IAM Config")
            .args([
                arg!(--"iam.api-version" "IAM api version, default is `v1`.")
                    .default_value("v1")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.server.timeout-ms" "")
                    .value_parser(clap::value_parser!(u16).range(10..))
                    .default_value("1000")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.server.max-retries" "")
                    .value_parser(clap::value_parser!(u8).range(0..10))
                    .default_value("3")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.server.retry-duration-ms" "")
                    .value_parser(clap::value_parser!(u16).range(10..))
                    .default_value("1000")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.server.address" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.server.use-tls" "")
                    .action(ArgAction::SetTrue)
                    .group(ID),
                arg!(--"iam.server.skip-tls-verify" "")
                    .action(ArgAction::SetTrue)
                    .group(ID),
                arg!(--"iam.server.tls-cert" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.server.tls-key" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.tls-cert" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.tls-cert-data" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.tls-key" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.tls-key-data" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.token" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.username" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.password" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.secret-id" "")
                    .action(ArgAction::Append)
                    .group(ID),
                arg!(--"iam.auth.secret-key" "")
                    .action(ArgAction::Append)
                    .group(ID),
            ])
    }

    /// Parse matches info `Config`.
    ///
    /// `command.get_matches()` need to get ownership of the command, so borrowing matches is ok.
    pub fn parse_command(matches: &ArgMatches) -> Self {
        Self {
            auth_info: AuthInfo {
                token: matches
                    .get_one::<String>("iam.auth.token")
                    .expect("must provide a token")
                    .clone(),
                username: matches
                    .get_one::<String>("iam.auth.username")
                    .map(|s| s.to_string()),
                password: matches
                    .get_one::<String>("iam.auth.password")
                    .map(|s| s.to_string()),
                secret_id: matches
                    .get_one::<String>("iam.auth.secret-id")
                    .map(|s| s.to_string()),
                secret_key: matches
                    .get_one::<String>("iam.auth.secret-key")
                    .map(|s| s.to_string()),
            },
            server_info: ServerInfo {
                timeout_ms: matches
                    .get_one::<u16>("iam.server.timeout-ms")
                    .unwrap()
                    .clone(),
                max_retries: matches
                    .get_one::<u8>("iam.server.max-retries")
                    .unwrap()
                    .clone(),
                retry_duration_ms: matches
                    .get_one::<u16>("iam.server.retry-duration-ms")
                    .unwrap()
                    .clone(),
                address: matches
                    .get_one::<String>("iam.server.address")
                    .unwrap()
                    .to_string(),
                use_tls: matches.get_flag("iam.server.use-tls"),
                skip_tls_verify: matches.get_flag("iam.server.skip-tls-verify"),
                tls_cert: matches
                    .get_one::<String>("iam.server.tls-cert")
                    .map(|s| s.to_string()),
                tls_cert_data: matches
                    .get_one::<String>("iam.server.tls-cert-data")
                    .map(|s| s.to_string()),
                tls_key: matches
                    .get_one::<String>("iam.server.tls-key")
                    .map(|s| s.to_string()),
                tls_key_data: matches
                    .get_one::<String>("iam.server.tls-key-data")
                    .map(|s| s.to_string()),
            },
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Clone)]
#[serde(rename_all = "kebab-case")]
pub struct AuthInfo {
    pub token: String,
    pub username: Option<String>,
    pub password: Option<String>,
    pub secret_id: Option<String>,
    pub secret_key: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
#[serde(rename_all = "kebab-case")]
pub struct ServerInfo {
    #[serde(default = "default_timeout_ms")]
    pub timeout_ms: u16,
    #[serde(default = "default_max_retries")]
    pub max_retries: u8,
    #[serde(default = "default_retry_duration_ms")]
    pub retry_duration_ms: u16,
    pub address: String,
    pub use_tls: bool,
    pub skip_tls_verify: bool,
    pub tls_cert: Option<String>,
    pub tls_cert_data: Option<String>,
    pub tls_key: Option<String>,
    pub tls_key_data: Option<String>,
}

fn default_timeout_ms() -> u16 {
    1000
}

fn default_max_retries() -> u8 {
    3
}

fn default_retry_duration_ms() -> u16 {
    1000
}

#[test]
fn test_config() {
    let c = Config::load_from_file("config.yaml");
    println!("{:?}", c);

    let cmd = Command::new("my_prog")
        .version("1.0")
        .author("Kevin K. <kbknapp@gmail.com>")
        .about("Does awesome things");

    let cmd = Config::add_to_command(cmd);
    let arg_vec = vec![
        "my_prog",
        "--iam.server.timeout-ms",
        "1000",
        "--iam.server.max-retries",
        "4",
        "--iam.server.retry-duration-ms",
        "2000",
        "--iam.server.address",
        "127.0.0.1",
        "--iam.server.use-tls",
        "--iam.server.skip-tls-verify",
        "--iam.server.tls-cert",
        "12312312",
        "--iam.server.tls-cert-data",
        "12312312",
        "--iam.server.tls-key",
        "12312312",
        "--iam.server.tls-key-data",
        "12312312",
        "--iam.auth.token",
        "12312312",
        "--iam.auth.username",
        "admin",
        "--iam.auth.password",
        "admin",
    ];
    let c = Config::parse_command(&cmd.get_matches_from(arg_vec));
    println!("{:?}", c);
}
