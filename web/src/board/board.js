export const symbolsMap = {
    2: ["marking", "32"],
    0: ["marking marking-x", 9587], // represents x
    1: ["marking marking-o", 9711], // represents o
};

export const patterns = [
    //horizontal
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    //vertical
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    //diagonal
    [0, 4, 8],
    [2, 4, 6]
];

export const AIScore = {2: 1, 0: 2, 1: 0};
