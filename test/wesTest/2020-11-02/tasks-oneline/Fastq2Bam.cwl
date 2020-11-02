cwlVersion: v1.0

class: CommandLineTool 
id: fastq_to_bam 

requirements:
  - class: DockerRequirement
    dockerPull: vinbdi/biotools
  - class: ResourceRequirement
    ramMin: 4096
    coresMin: 4

inputs:
  input_fastqs1:
    type: File[]
  input_fastqs2:
    type: File[]
  output_basename:
    type: string
  output_fastq2bam_bam:
    type: string

outputs:
  output: 
    type: File 
    outputBinding:
      glob: '*unmapped.bam'
      
arguments: 
  - position: 0
    shellQuote: false 
    valueFrom: >- 
      SM=$(inputs.output_basename) RG=$(inputs.output_basename)  LIBRARY_NAME=lib1 PLATFORM=ILLUMINA F1=$(inputs.input_fastqs1) F2=$(input.input_fastqs2) O=$(input.output_fastq2bam_bam)
baseCommand: [picard, FastqToSam]
