import argparse 

def usage() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Compare operator events")
    parser.add_argument("success_dir")
    parser.add_argument("failed_dir")
    parser.add_argument("--threshold", default=0.2)
    return parser.parse_args()

args1 = usage()
print(args1.echo)

