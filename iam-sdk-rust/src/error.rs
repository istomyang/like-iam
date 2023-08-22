use std::{
    error,
    fmt::{self, Display},
    result,
};

/// A specialized [`Result`] type for iam-sdk-rust operations.
///
/// This type is broadly used across [`crate`] for any operation which may
/// produce an error.
pub type Result<T> = result::Result<T, Error>;

#[derive(Debug, Clone)]
pub struct Error {
    kind: ErrorKind,
    father: Option<Box<Error>>,
    cause: Option<String>, // If use error::Error directly, ownership is a little complex.
}

impl Error {
    pub fn new(kind: ErrorKind, cause: Option<Box<dyn error::Error>>) -> Self {
        Error {
            kind,
            father: None,
            cause: cause.map(|e| e.to_string()),
        }
    }

    pub fn new_with_str(kind: ErrorKind, cause: String) -> Self {
        Error {
            kind,
            father: None,
            cause: Some(cause),
        }
    }

    pub fn wrap(&self, kind: ErrorKind, cause: Option<Box<dyn error::Error>>) -> Self {
        Error {
            kind,
            father: Some(Box::new(self.clone())),
            cause: cause.map(|e| e.to_string()),
        }
    }

    pub fn is(&self, kind: ErrorKind) -> bool {
        let mut cur = self;
        loop {
            if cur.kind == kind {
                return true;
            }
            if cur.father.is_some() {
                cur = cur.father.as_ref().unwrap().as_ref();
                continue;
            }
            break;
        }
        false
    }

    pub fn cause_string(&self) -> Option<&String> {
        self.cause.as_ref()
    }
}

impl From<ErrorKind> for Error {
    fn from(kid: ErrorKind) -> Self {
        Self::new(kid, None)
    }
}

impl Default for Error {
    fn default() -> Self {
        Self {
            kind: ErrorKind::Unknown,
            father: None,
            cause: None,
        }
    }
}

impl error::Error for Error {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        Some(self)
    }

    fn description(&self) -> &str {
        self.kind.as_str()
    }

    fn cause(&self) -> Option<&dyn error::Error> {
        self.source()
    }
}

impl Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.kind.as_str())?;

        let mut cur = &self.father;
        while let Some(e) = cur {
            f.write_str(";")?;
            f.write_str(e.kind.as_str())?;
            cur = &e.father;
        }
        f.write_str(";")?;
        Ok(())
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, Ord, PartialEq, PartialOrd)]
#[non_exhaustive]
pub enum ErrorKind {
    JsonDecode,
    JsonEncode,
    OpenFile,
    ParseError,
    InvalidVersion,
    Unknown,
    Network,
    _NonExhaustive,
}

impl ErrorKind {
    pub fn as_str(&self) -> &'static str {
        match *self {
            ErrorKind::_NonExhaustive => "_NonExhaustive.",
            ErrorKind::JsonDecode => "json decode error.",
            ErrorKind::JsonEncode => "json encode error.",
            ErrorKind::OpenFile => "open file error.",
            ErrorKind::ParseError => "parse error.",
            ErrorKind::InvalidVersion => "invalid version error.",
            ErrorKind::Unknown => "unknown error.",
            ErrorKind::Network => "network error.",
        }
    }
}

impl fmt::Display for ErrorKind {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.as_str())
    }
}
