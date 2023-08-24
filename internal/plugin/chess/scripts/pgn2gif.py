#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
import pgn2gif


def generate_gif(pgn: str, out_dir: str, file_name: str):
    creator = pgn2gif.PgnToGifCreator(
        reverse=False, duration=1, ws_color='white', bs_color='gray')
    pgn_file_path = os.path.join(out_dir, file_name + ".pgn")
    with open(pgn_file_path, "w") as f:
        f.write(pgn)
    out_gif_path = os.path.join(out_dir, file_name + ".gif")
    creator.create_gif(pgn_file_path, out_path=out_gif_path)


if __name__ == "__main__":
    args = sys.argv
    if len(args) != 4:
        print("Wrong args")
        exit(1)
    pgn_str = args[1]
    out_dir = args[2]
    file_name = args[3]
    generate_gif(pgn_str, out_dir, file_name)
    exit(0)
