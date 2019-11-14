import socket
import select
import sys
import threading


HEADER_SIZE = 10
PORT = 1234

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

server.bind((socket.gethostname(), PORT))
server.listen(10)

clients = {}


def handle_client(client):
    # msg = f'Welcome to server. Your IP is {clients[client]["addr"]}'
    while True:
        read_sockets, _, _ = select.select([sys.stdin, client], [], [])
        for sock in read_sockets:
            if sock == client:
                full_msg = ''
                init_msg_received = False

                while True:
                    msg = client.recv(16).decode('utf-8')

                    if not init_msg_received:
                        msg_len = int(msg[:HEADER_SIZE])
                        init_msg_received = True

                    full_msg += msg

                    if len(full_msg) - HEADER_SIZE == msg_len:
                        print(f'Received: {full_msg}')
                        break
            else:
                msg = input()
                msg = f'{len(msg):<{HEADER_SIZE}}{msg}'
                for conn in clients.values():
                    conn['conn'].send(msg.encode('utf-8'))


while True:
    conn, addr = server.accept()
    print(f'Connected to: {addr}')

    clients[addr[1]] = {}
    clients[addr[1]]['conn'] = conn
    thread = threading.Thread(target=handle_client, args=(conn, ))
    thread.start()

    clients[addr[1]]['thread'] = thread
    clients[addr[1]]['addr'] = addr
