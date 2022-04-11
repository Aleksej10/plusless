mod config;

use rblocks::*;
use std::process;

fn main() {
  let (blocks, signals) = initialize().unwrap_or_else(|err| {
    eprintln!("{}", err);
    process::exit(1);
  });

  blocks.run();
  blocks.draw();
  blocks.listen(signals);
}
