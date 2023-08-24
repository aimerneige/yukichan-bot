#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import chess
import chess.svg
import datetime

gray_theme = {
    "square light": "#F5F5F5",
    "square dark": "#BDBDBD",
    "square light lastmove": "#616161",
    "square dark lastmove": "#5C5C5C",
    "margin": "#212121",
    "coord": "#E2E2E2",
}
green_theme = {
    "square light": "#769656",
    "square dark": "#EEEED2",
    "square light lastmove": "#BACA2B",
    "square dark lastmove": "#F6F669",
    "margin": "#212121",
    "coord": "#E2E2E2",
}
red_theme = {
    "square light": "#f2cfb6",
    "square dark": "#c24539",
    "square light lastmove": "#e16b8c",
    "square dark lastmove": "#f06c91",
    "margin": "#212121",
    "coord": "#E2E2E2",
}


def generate_board_svg(fen_str: str, file_path: str, last_move_uci: str):
    board = chess.Board(fen_str)

    king_square = None
    if not is_12_13_day() and board.is_check():
        king_square_index = board.king(board.turn)
        king_square_name = chess.square_name(king_square_index)
        king_square = chess.parse_square(king_square_name)

    last_move = None
    if last_move_uci != "None":
        last_move = chess.Move.from_uci(last_move_uci)

    themes = {}
    if is_christmas_day() or is_new_year_day():
        themes = red_theme
    elif is_april_fools_day():
        themes = green_theme
    elif is_12_13_day():
        themes = gray_theme

    svg = chess.svg.board(
        board=board,
        orientation=board.turn,
        lastmove=last_move,
        check=king_square,
        size=720,
        coordinates=True,
        colors=themes,
    )
    with open(file_path, 'w') as f:
        f.write(svg)


def is_christmas_day() -> bool:
    now = datetime.datetime.now()
    return now.month == 12 and now.day == 25


def is_new_year_day() -> bool:
    now = datetime.datetime.now()
    return now.month == 1 and now.day == 1


def is_april_fools_day() -> bool:
    now = datetime.datetime.now()
    return now.month == 4 and now.day == 1


def is_12_13_day() -> bool:
    now = datetime.datetime.now()
    return now.month == 12 and now.day == 13


if __name__ == "__main__":
    args = sys.argv
    if len(args) != 4:
        print("Wrong args")
        exit(1)
    fen_str = args[1]
    file_path = args[2]
    last_move_uci = args[3]
    generate_board_svg(fen_str, file_path, last_move_uci)
    exit(0)
