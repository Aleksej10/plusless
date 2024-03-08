mod config;
use crate::config::*;
use signal_hook::iterator::Signals;
use std::process::Command;
use std::sync::{Arc, RwLock};
use std::thread;
use std::time::Duration;
use x11;

const SIGRTMIN: i32 = 34;

pub struct CfgBlk {
    pub sig: i32,
    pub t: i32,
    pub ico: &'static str,
    pub cmd: &'static str,
}

#[derive(Clone)]
struct Blk {
    sig: i32,
    t: i32,
    ico: String,
    cmd: &'static str,
}

struct Sig {
    blk: i32,
    btn: Option<i32>,
}

struct Block {
    blk: Blk,
    res: Arc<RwLock<String>>,
    should_update: Arc<RwLock<bool>>,
}

impl Block {
    fn new(blk: Blk, should_update: Arc<RwLock<bool>>) -> Block {
        Block {
            blk,
            res: Arc::new(RwLock::new(String::from(""))),
            should_update,
        }
    }

    fn to_s(&self) -> Option<String> {
        match self.res.read() {
            Ok(res) => {
                if res.is_empty() {
                    None
                } else {
                    Some(format!("{}{}", self.blk.ico, res))
                }
            }
            Err(e) => {
                let res = e.into_inner().to_owned();
                if res.is_empty() {
                    None
                } else {
                    Some(format!("{}{}TAINTED", self.blk.ico, res)) // TODO: what do?
                }
            }
        }
    }

    fn set_should_update(&self, new_su: bool) {
        if let Ok(su) = self.should_update.read() {
            if *su == new_su {
                return;
            }
        }
        if let Ok(mut su) = self.should_update.write() {
            *su = new_su;
        }
    }

    fn update_res(&self, new_res: String) {
        if let Ok(old_res) = self.res.read() {
            if *old_res == new_res {
                return;
            }
        }
        if let Ok(mut old_res) = self.res.write() {
            *old_res = new_res;
        }
        self.set_should_update(true);
    }

    fn clone(&self) -> Block {
        Block {
            blk: self.blk.clone(),
            res: Arc::clone(&self.res),
            should_update: Arc::clone(&self.should_update),
        }
    }

    fn exec(&self, btn_sig: Option<i32>) {
        self.update_res(exec_cmd(self.blk.cmd, btn_sig));
    }

    fn run(&self) -> thread::JoinHandle<()> {
        let blk = self.clone();

        thread::spawn(move || {
            blk.update_res(exec_cmd(blk.blk.cmd, None));
            if blk.blk.t == 0 {
                return;
            }

            loop {
                thread::sleep(Duration::from_secs(blk.blk.t as u64));
                blk.update_res(exec_cmd(blk.blk.cmd, None));
            }
        })
    }
}

pub struct Blocks {
    data: Vec<Block>,
    should_update: Arc<RwLock<bool>>,
}

impl Blocks {
    fn new(blocks: Vec<Blk>) -> Blocks {
        let should_update = Arc::new(RwLock::new(false));

        let data = blocks
            .into_iter()
            .map(|blk| Block::new(blk, Arc::clone(&should_update)))
            .collect();

        Blocks {
            data,
            should_update,
        }
    }

    fn set_should_update(&self, new_su: bool) {
        if let Ok(su) = self.should_update.read() {
            if *su == new_su {
                return;
            }
        }
        if let Ok(mut su) = self.should_update.write() {
            *su = new_su;
        }
    }

    fn to_s(&self) -> String {
        self.set_should_update(false);

        self.data
            .iter()
            .map(|e| e.to_s())
            .filter(|e| e.is_some())
            .map(|e| e.unwrap())
            .collect::<Vec<String>>()
            .join(DELIM)
    }

    fn clone(&self) -> Blocks {
        let data = self.data.iter().map(|b| b.clone()).collect();
        let should_update = Arc::clone(&self.should_update);

        Blocks {
            data,
            should_update,
        }
    }

    fn exec_block(&self, sig: Sig) -> Option<thread::JoinHandle<()>> {
        let blk = self.data.iter().find(|&block| block.blk.sig == sig.blk)?;
        let blk = blk.clone();
        let blocks = self.clone();

        Some(thread::spawn(move || {
            blk.exec(sig.btn);
            x11_set_status(blocks.to_s());
        }))
    }

    pub fn run(&self) -> Vec<thread::JoinHandle<()>> {
        self.data.iter().map(|blk| blk.run()).collect()
    }

    pub fn draw(&self) -> thread::JoinHandle<()> {
        let blocks = self.clone();

        thread::spawn(move || loop {
            thread::sleep(Duration::from_secs(2));

            if let Ok(su) = blocks.should_update.read() {
                if *su == false {
                    continue;
                }
            }
            x11_set_status(blocks.to_s());
        })
    }

    pub fn listen(&self, mut signals: Signals) {
        let mut iter = signals.into_iter();

        loop {
            if let Some(sig) = parse_signal(&mut iter) {
                self.exec_block(sig);
            }
        }
    }
}

pub fn initialize() -> Result<(Blocks, Signals), &'static str> {
    let mut blocks = Vec::<Blk>::with_capacity(30);
    let mut tmp_sig: i32 = 6;
    let mut signals = vec![35, 36, 37, 38, 39, 40];

    for blk in BLOCKS.iter() {
        let sig = if blk.sig != 0 {
            if blk.sig < 7 {
                return Err("use signals in range [7..30]");
            }
            blk.sig
        } else {
            tmp_sig = next_tmp_sig(tmp_sig)?;
            tmp_sig
        };

        let byte = String::from_utf8_lossy(&[sig as u8]).to_string();
        let ico = if blk.ico.is_empty() {
            byte
        } else {
            format!("{}{} ", byte, blk.ico.trim())
        };

        blocks.push(Blk {
            sig,
            t: blk.t,
            ico,
            cmd: blk.cmd,
        });
        signals.push(SIGRTMIN + sig);
    }

    let signals = Signals::new(&signals).map_err(|_| "error while binding signals")?;
    let blocks = Blocks::new(blocks);
    return Ok((blocks, signals));
}

fn next_tmp_sig(tmp_sig: i32) -> Result<i32, &'static str> {
    let mut next_sig = tmp_sig + 1;

    while next_sig < 31 {
        let mut ok = true;

        for blk in BLOCKS.iter() {
            if next_sig == blk.sig {
                next_sig += 1;
                ok = false;
                break;
            }
        }
        if ok {
            return Ok(next_sig);
        }
    }
    Err("you probably have waaay too many blocks")
}

fn exec_cmd(cmd: &'static str, sig: Option<i32>) -> String {
    let mut command = Command::new(SHELL);
    command.arg("-c");
    command.arg(cmd);

    if let Some(s) = sig {
        command.env("BLOCK_BUTTON", s.to_string());
    }
    match command.output() {
        Ok(o) => String::from_utf8_lossy(&o.stdout).trim().to_string(),
        Err(_) => String::from(""),
    }
}

fn x11_set_status(text: String) {
    if let Ok(c_status) = std::ffi::CString::new(text) {
        unsafe {
            let d = x11::xlib::XOpenDisplay(std::ptr::null());
            let screen = x11::xlib::XDefaultScreenOfDisplay(d);
            let root = x11::xlib::XRootWindowOfScreen(screen);
            x11::xlib::XStoreName(d, root, c_status.into_raw());
            x11::xlib::XCloseDisplay(d);
        }
    }
}

fn parse_signal(
    iter: &mut signal_hook::iterator::Forever<signal_hook::iterator::exfiltrator::SignalOnly>,
) -> Option<Sig> {
    let blk = iter.next()? - SIGRTMIN;
    if blk < 7 {
        return None;
    }

    let btn = iter.next()? - SIGRTMIN;

    if blk == btn {
        return Some(Sig { blk, btn: None });
    }
    if btn < 7 {
        return Some(Sig {
            blk,
            btn: Some(btn),
        });
    }
    None
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn exec_cmd_test_trim_test() {
        let cmd = "echo '     '";
        let out = exec_cmd(cmd, None);

        assert_eq!(out, String::from(""));
    }

    #[test]
    fn exec_cmd_test_error() {
        let cmd = "cyka blyat";
        let out = exec_cmd(cmd, None);

        assert_eq!(out, String::from(""));
    }

    #[test]
    fn exec_cmd_test_button() {
        for i in 1..7 {
            let cmd = "echo $BLOCK_BUTTON";
            let out = exec_cmd(cmd, Some(i));

            assert_eq!(out, i.to_string());
        }
    }

    #[test]
    fn x11_set_status_test() {
        x11_set_status(String::from("herro!"));
    }
}
