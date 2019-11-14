import socket
import sys
import select

client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

while True:
    read_streams, _, _ = select.select([sys.stdin, client], [], [])

    for stream in read_streams:
        if stream == client:
            msg, addr = client.recvfrom(4096)
            msg = msg.decode('utf-8')
            print(f'Received: {msg}')
        else:
            msg = input()
            msg = msg.encode('utf-8')
            client.sendto(msg, (socket.gethostname(), 1234))