import socket
import sys
import select

server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
server.bind((socket.gethostname(), 1234))

clients = {}

while True:
    read_streams, _, _ = select.select([sys.stdin, server], [], [])
    for stream in read_streams:
        if stream == server:
            data, addr = server.recvfrom(4096)
            if clients.get(addr[1]) is None:
                print(f'Connected to: {addr}')
                clients[addr[1]] = addr

            data = data.decode('utf-8')

            for address in clients.values():
                server.sendto(data.encode('utf-8'), address)
        else:
            msg = input()
            msg = msg.encode('utf-8')
            
            for address in clients.values():
                server.sendto(data, address)
