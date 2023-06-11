# import re
# from shapely import wkb
# import binascii
#
# def convert_hex_to_point(hex_str):
#     bytes_wkb = binascii.unhexlify(hex_str)
#     point = wkb.loads(bytes_wkb)
#     return "ST_GeomFromText('{}')".format(point.wkt)
#
# with open('./spots.sql', 'r') as f:
#     lines = f.readlines()
#
# with open('./newspots.sql', 'w') as f:
#     for line in lines:
#         match = re.search(r'([0-9a-fA-F]{50,})', line)
#         if match:
#             new_line = line.replace("'{}'".format(match.group()), convert_hex_to_point(match.group()))
#             f.write(new_line)
#         else:
#             f.write(line)


def count_insert_queries(file_path):
    with open(file_path, 'r') as file:
        content = file.read()
        queries = content.split("INSERT INTO")
        return len(queries) - 1

print(count_insert_queries('./newspots.sql'))
