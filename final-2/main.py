import sys
from ortools.linear_solver import pywraplp


def read_input_file(input_path):
    with open(input_path, "r") as file:
        n = int(file.readline())
        m = int(file.readline())
        edges = [tuple(map(int, line.strip().split())) for line in file]
    return n, m, edges


def write_output_file(output_path, k, s):
    with open(output_path, "w") as file:
        file.write(f"{k}:{s}")


def minimum_dominating_set(input_path, output_path):
    n, m, edges = read_input_file(input_path)
    solver = pywraplp.Solver.CreateSolver('SCIP')
    nodes = [solver.BoolVar(str(i)) for i in range(n)]
    adjacency_list = {i: [] for i in range(n)}

    for a, b in edges:
        adjacency_list[a].append(b)
        adjacency_list[b].append(a)

    for i in range(n):
        x = solver.IntVar(1, n, f'x{i}')
        constraint = solver.Constraint(0, 0, f'ct{i}')
        constraint.SetCoefficient(x, -1)
        constraint.SetCoefficient(nodes[i], 1)

        for neighbor in adjacency_list[i]:
            constraint.SetCoefficient(nodes[neighbor], 1)

    objective = solver.Objective()
    for node in nodes:
        objective.SetCoefficient(node, 1)
    objective.SetMinimization()

    solver.Solve()

    k = int(objective.Value())
    s = ''.join(str(int(node.solution_value())) for node in nodes)

    write_output_file(output_path, k, s)


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("usage: prog [inputpath] [outputpath]")
    else:
        input_path, output_path = sys.argv[1], sys.argv[2]
        minimum_dominating_set(input_path, output_path)
