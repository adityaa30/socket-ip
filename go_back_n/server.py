import socket

ADDR = (socket.gethostname(), 1236)

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(ADDR)
server.listen(5)

WINDOW_SIZE = 4
expect_val = 0
send_msg = ''

try:
    conn, addr = server.accept()
    print "Connected to " + str(addr)

    while True:
        msg = conn.recv(4096).decode('utf-8')
        orig_msg = msg
        ok = True
        for num in msg.split(','):
            if int(num) == expect_val:
                expect_val = (expect_val + 1) % WINDOW_SIZE
                send_msg = 'ACK ' + str(expect_val)                
            else:
                send_msg = 'NCK ' + str(expect_val)
                break
        
        conn.send(send_msg.encode('utf-8'))

except Exception as e:
    print e
    conn.close()
    