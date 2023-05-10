FROM busybox

COPY build/toy-lustre-csi /

ENTRYPOINT ["/toy-lustre-csi"]