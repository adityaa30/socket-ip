import socket

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind((socket.gethostname(), 1234))

while True:
    data, addr = sock.recvfrom(4096)

    print(f'Received \'{data.decode('utf-8')}\' from \'{addr}\'')

    if len(data) > 0:
        sent = sock.sendto(data, addr)

