FROM wealdtech/ethdo:1.36.1 as ethdo

FROM attestant/vouch:1.9.2

COPY --from=ethdo /app/ethdo /app/ethdo

RUN apt-get update && apt-get install -y curl jq wget

ENTRYPOINT ["/compose/vouch/run.sh"]
